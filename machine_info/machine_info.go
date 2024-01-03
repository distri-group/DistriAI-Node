package machine_info

import (
	"DistriAI-Node/machine_info/cpu"
	"DistriAI-Node/machine_info/disk"
	"DistriAI-Node/machine_info/flops"
	"DistriAI-Node/machine_info/gpu"
	"DistriAI-Node/machine_info/location"
	"DistriAI-Node/machine_info/machine_uuid"
	"DistriAI-Node/machine_info/memory"
	"DistriAI-Node/machine_info/speedtest"
	logs "DistriAI-Node/utils/log_utils"
)

// MachineInfo 用于存储所有硬件(cpu\disk\flops\gpu\memory)信息
type MachineInfo struct {
	MachineUUID     machine_uuid.MachineUUID `json:"MachineUUID"`     // 机器 UUID
	Addr            string                   `json:"Addr"`            // 用户钱包地址
	CPUInfo         cpu.InfoCPU              `json:"CPUInfo"`         // CPU 信息
	DiskInfo        disk.InfoDisk            `json:"DiskInfo"`        // 硬盘信息
	Score           float64                  `json:"Score"`           // 得分
	MemoryInfo      memory.InfoMemory        `json:"InfoMemory"`      // 内存信息
	GPUInfo         gpu.InfoGPU              `json:"GPUInfo"`         // GPU 信息（仅限英特尔显卡）
	LocationInfo    location.InfoLocation    `json:"LocationInfo"`    // IP 对应的地理位置
	SpeedInfo       speedtest.InfoSpeed      `json:"SpeedInfo"`       // 上传下载速度
	FlopsInfo       flops.InfoFlop           `json:"InfoFlop"`        // FLOPS 信息
	MachineAccounts string                   `json:"MachineAccounts"` // 机器存储地址
}

// GetMachineInfo 函数收集并返回全部硬件信息
func GetMachineInfo() (MachineInfo, error) {
	var hwInfo MachineInfo

	// 获取 CPU 信息
	cpuInfo, err := cpu.GetCPUInfo()
	if err != nil {
		return hwInfo, err
	}
	hwInfo.CPUInfo = cpuInfo

	// 获取内存信息
	memInfo, err := memory.GetMemoryInfo()
	if err != nil {
		return hwInfo, err
	}
	hwInfo.MemoryInfo = memInfo

	// 获取 GPU 信息
	gpuInfo, _ := gpu.GetIntelGPUInfo()
	// if err != nil {
	// 	return hwInfo, err
	// }
	hwInfo.GPUInfo = gpuInfo

	// 获取 IP 对应的地理位置信息
	locationInfo, err := location.GetLocationInfo()
	if err != nil {
		return hwInfo, err
	}
	hwInfo.LocationInfo = locationInfo

	// Easy debugging
	// 测算网络上传下载速率
	speedInfo, err := speedtest.GetSpeedInfo()
	if err != nil {
		logs.Error(err.Error())	
	}
	// speedInfo := speedtest.InfoSpeed{
	// 	Download: "1000 Mbit/s",
	// 	Upload:   "1000 Mbit/s",
	// }
	hwInfo.SpeedInfo = speedInfo

	// 获取 FLOPS 信息

	flopsInfo := flops.GetFlopsInfo(gpuInfo.Model)
	hwInfo.FlopsInfo = flopsInfo

	return hwInfo, nil
}
