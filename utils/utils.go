package utils

import (
	docker_utils "DistriAI-Node/docker/utils"
	"DistriAI-Node/machine_info/machine_uuid"
	"DistriAI-Node/pattern"
	logs "DistriAI-Node/utils/log_utils"
	"archive/zip"
	"crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/cavaliergopher/grab/v3"
	"github.com/cavaliergopher/grab/v3/pkg/grabtest"
)

func ByteUUIDToStrUUID(byteUUID pattern.MachineUUID) machine_uuid.MachineUUID {
	return machine_uuid.MachineUUID(hex.EncodeToString(byteUUID[:]))
}

func ParseMachineUUID(uuidStr string) (pattern.MachineUUID, error) {
	/* Linux */
	var machineUUID pattern.MachineUUID

	b, err := hex.DecodeString(uuidStr)
	if err != nil {
		return machineUUID, fmt.Errorf("> hex.DecodeString: %v", err.Error())
	}
	copy(machineUUID[:], b[:16])

	return machineUUID, nil
}

func ParseTaskUUID(uuidStr string) (pattern.TaskUUID, error) {
	/* Linux */
	var taskUUID pattern.TaskUUID

	b, err := hex.DecodeString(uuidStr)
	if err != nil {
		panic(err)
	}
	copy(taskUUID[:], b[:16])

	return taskUUID, nil
}

func Zip(src, dest string) error {
	destFile, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer func(destFile *os.File) {
		err := destFile.Close()
		if err != nil {

		}
	}(destFile)

	myZip := zip.NewWriter(destFile)
	defer func(myZip *zip.Writer) {
		err := myZip.Close()
		if err != nil {

		}
	}(myZip)

	err = filepath.Walk(src, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		relPath, err := filepath.Rel(src, filePath)
		if err != nil {
			return err
		}

		zipFile, err := myZip.Create(relPath)
		if err != nil {
			return err
		}

		fsFile, err := os.Open(filePath)
		if err != nil {
			return err
		}
		defer func(fsFile *os.File) {
			err := fsFile.Close()
			if err != nil {

			}
		}(fsFile)

		_, err = io.Copy(zipFile, fsFile)
		return err
	})
	if err != nil {
		return err
	}

	return nil
}

func EnsureHttps(url string) string {
	if !strings.HasPrefix(url, "https://") {
		return "https://" + url
	}
	return url
}

func Unzip(src string, dest string) ([]string, error) {
	var filenames []string

	r, err := zip.OpenReader(src)
	if err != nil {
		return filenames, err
	}
	defer func(r *zip.ReadCloser) {
		err := r.Close()
		if err != nil {
			return
		}
	}(r)

	for _, f := range r.File {
		fpath := filepath.Join(dest, f.Name)

		if !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
			return filenames, fmt.Errorf("illegal file path: %s", fpath)
		}

		filenames = append(filenames, fpath)

		if f.FileInfo().IsDir() {
			err := os.MkdirAll(fpath, os.ModePerm)
			if err != nil {
				return nil, err
			}
			continue
		}

		if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return filenames, err
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return filenames, err
		}

		rc, err := f.Open()
		if err != nil {
			return filenames, err
		}

		_, err = io.Copy(outFile, rc)
		if err != nil {
			return nil, err
		}

		err = outFile.Close()
		if err != nil {
			return nil, err
		}
		err = rc.Close()
		if err != nil {
			return nil, err
		}
	}
	return filenames, nil
}

func GetFreeSpace(path string) (uint64, error) {
	var stat syscall.Statfs_t
	err := syscall.Statfs(path, &stat)
	if err != nil {
		return 0, err
	}
	return stat.Bavail * uint64(stat.Bsize), nil
}

func CheckPort(port string) bool {
	logs.Normal(fmt.Sprintf("Checking port %s...", port))

	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return false
	}
	listener.Close()
	return true
}

func GenerateRandomString(length int) (string, error) {
	bytes := make([]byte, length/2)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", fmt.Errorf("> rand.Read: %v", err.Error())
	}
	return hex.EncodeToString(bytes), nil
}

func CompareSpaceWithDocker(sizeLimitGB int) (bool, error) {
	dockerSizeStr, err := docker_utils.GetDockerImageDirSize()
	if err != nil {
		return false, err
	}

	sizeLimitBytes := int64(sizeLimitGB) * 1024 * 1024 * 1024

	dockerSizeStr = strings.TrimSuffix(dockerSizeStr, "G")
	dockerSize, err := strconv.ParseFloat(dockerSizeStr, 64)
	if err != nil {
		return false, err
	}

	if int64(dockerSize*1024*1024*1024) < sizeLimitBytes {
		return false, nil
	}

	return true, nil
}

const (
	genesisTime    int64 = 1708992000
	periodDuration int64 = 86400
)

func CurrentPeriod() uint32 {
	return uint32((time.Now().Unix() - genesisTime) / periodDuration)
}

func PeriodBytes() []byte {
	bytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(bytes, CurrentPeriod())
	return bytes
}

func GetFilenameFromURL(rawURL string) (string, error) {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}
	return path.Base(parsedURL.Path), nil
}

func SplitURL(rawURL string) (string, string, error) {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return "", "", err
	}
	host := parsedURL.Scheme + "://" + parsedURL.Host
	path := parsedURL.Path
	return host, path, nil
}

type DownloadURL struct {
	URL      string
	Checksum string
}

func DownloadFiles(dest string, urls []DownloadURL) error {
	client := grab.NewClient()
	reqs := make([]*grab.Request, len(urls))

	for i, url := range urls {
		label, err := GetFilenameFromURL(url.URL)
		if err != nil {
			return err
		}

		req, err := grab.NewRequest(dest+"/"+label, url.URL)
		if err != nil {
			return err
		}

		req.Label = label
		// req.SetChecksum(sha256.New(), grabtest.MustHexDecodeString(url.Checksum), true)
		req.SetChecksum(nil, grabtest.MustHexDecodeString(url.Checksum), true)
		reqs[i] = req
	}

	responses := client.DoBatch(len(reqs), reqs...)

	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	var completed int
	for i := 0; i < len(reqs); {
		select {
		case resp := <-responses:
			if resp == nil {
				return fmt.Errorf("> resp is nil")
			}

			if err := resp.Err(); err != nil {
				return fmt.Errorf("> %s resp.Err: %v", resp.Request.Label, err.Error())
			}

			logs.Normal(fmt.Sprintf("%s (%.2f%%)", resp.Request.Label, 100*resp.Progress()))

			if resp.IsComplete() {
				completed++
			}
			if completed == len(reqs) {
				logs.Normal("All downloads completed")
				return nil
			}
			i++
		}
	}
	return errors.New("> DownloadFiles: unexpected exit")
}
