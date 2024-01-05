package disk

import (
	"DistriAI-Node/config"
	"DistriAI-Node/utils"
	logs "DistriAI-Node/utils/log_utils"
	"fmt"
	"os"
)

type InfoDisk struct {
	Path       string  `json:"Path"`
	TotalSpace float64 `json:"TotalSpace"`
	// FreeSpace  float64 `json:"FreeSpace"`
}

func GetDiskInfo() (InfoDisk, error) {
	logs.Normal("Getting free space info...")

	dirpath := config.GlobalConfig.Console.WorkDirectory
	if _, err := os.Stat(dirpath); os.IsNotExist(err) {
		logs.Normal(fmt.Sprintf("%s does not exist. Using default directory /home", dirpath))
		dirpath = "/home"
	}

	freeSpace, err := utils.GetFreeSpace(dirpath)
	if err != nil {
		return InfoDisk{}, fmt.Errorf("error calculating free space: %v", err)
	}

	diskInfo := InfoDisk{
		Path:       dirpath,
		TotalSpace: float64(freeSpace) / 1024 / 1024 / 1024,
	}
	return diskInfo, nil
}
