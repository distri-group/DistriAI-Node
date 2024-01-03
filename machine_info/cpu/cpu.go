package cpu

import (
	"DistriAI-Node/utils/log_utils"
	"fmt"

	"github.com/shirou/gopsutil/v3/cpu"
)

type InfoCPU struct {
	ModelName string  `json:"ModelName"`
	Cores     int32   `json:"Cores"`
	Mhz       float64 `json:"Mhz"`
}

func GetCPUInfo() (InfoCPU, error) {
	logs.Normal("Getting CPU info...")

	cpuInfoStats, err := cpu.Info()
	if err != nil {
		return InfoCPU{}, fmt.Errorf("failed to get CPU info: %w", err)
	}

	counts, err := cpu.Counts(false)
	if err != nil {
		return InfoCPU{}, fmt.Errorf("failed to get CPU Counts : %w", err)
	}

	// cpuInfos := make([]InfoCPU, 0)
	// for _, info := range cpuInfoStats {
	// 	cpuInfos = append(cpuInfos, InfoCPU{
	// 		ModelName: info.ModelName,
	// 		Cores:     info.Cores,
	// 		Mhz:       info.Mhz,
	// 	})
	// }

	cpuInfo := InfoCPU{
		ModelName: cpuInfoStats[0].ModelName,
		Cores:     int32(counts),
		Mhz:       cpuInfoStats[0].Mhz,
	}

	return cpuInfo, nil
}
