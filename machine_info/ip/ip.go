package ip

import (
	"DistriAI-Node/config"
	logs "DistriAI-Node/utils/log_utils"
	"fmt"
	"os/exec"
	"strings"
)

type InfoIP struct {
	IP   string `json:"ip"`
	Port string `json:"port"`
}

func GetIpInfo() (InfoIP, error) {
	logs.Normal("Getting outer net ip info...")

	cmd := exec.Command("curl", "cip.cc")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return InfoIP{}, fmt.Errorf("error getting ip info: %v", err)
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "IP") {
			parts := strings.Split(line, ":")
			if len(parts) > 1 {
				ip := strings.TrimSpace(parts[1])
				return InfoIP{
					IP:   ip,
					Port: config.GlobalConfig.Console.OuterNetPort,
				}, nil
			}
		}
	}
	return InfoIP{}, fmt.Errorf("no outer net IP found")
}
