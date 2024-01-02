package speedtest

import (
	"DistriAI-Node/utils/log_utils"
	"fmt"
	"os/exec"
	"strings"
)

// InfoSpeed 定义 InfoSpeed 结构体
type InfoSpeed struct {
	Download string `json:"Download"` // 下载速度
	Upload   string `json:"Upload"`   // 上次速度
}

// GetSpeedInfo 测试网络的上传下载速率并返回 InfoSpeed 结构体
func GetSpeedInfo() (InfoSpeed, error) {
	logs.Normal("Getting network speed info...")

	out, err := exec.Command("speedtest-cli").Output()

	if err != nil {
		return InfoSpeed{
			Download: "- Mbit/s",
			Upload:   "- Mbit/s",
		}, fmt.Errorf("speedtest-cli, Failed to execute command: %v\nout: %v", err, string(out))
	}

	output := string(out)
	lines := strings.Split(output, "\n")

	var netSpeed InfoSpeed

	for _, line := range lines {
		if strings.Contains(line, "Download") {
			netSpeed.Download = strings.TrimSpace(strings.TrimPrefix(line, "Download:"))
		}

		if strings.Contains(line, "Upload") {
			netSpeed.Upload = strings.TrimSpace(strings.TrimPrefix(line, "Upload:"))
		}
	}
	return netSpeed, nil
}
