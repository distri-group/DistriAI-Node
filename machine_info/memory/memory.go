package memory

import (
	"DistriAI-Node/utils/log_utils"
	"fmt"

	"github.com/shirou/gopsutil/v3/mem"
)

// InfoMemory 定义 InfoMemory 结构体
type InfoMemory struct {
	RAM float64 `json:"RAM"` // 总内存大小，单位 GB
}

// GetMemoryInfo 获取内存信息并返回 InfoMemory 结构体
func GetMemoryInfo() (InfoMemory, error) {
	logs.Normal("Getting memory info...")

	vmStat, err := mem.VirtualMemory()
	if err != nil {
		return InfoMemory{}, fmt.Errorf("failed to get memory info: %w", err)
	}

	return InfoMemory{
		RAM: float64(vmStat.Total) / (1 << 30),
	}, nil
}
