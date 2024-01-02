package flops

import (
	"DistriAI-Node/machine_info/gpu/gpu_infos"
	"DistriAI-Node/utils/log_utils"
	"strings"
)

// InfoFlop 结构体用于存储 FLOPS 计算结果
type InfoFlop struct {
	Flops float32 `json:"Flops"` // 每秒浮点运算次数
}

// GetFlopsInfo 函数计算 FLOPS 并返回包含结果的 InfoFlop 结构体，接受 CPU 核心数作为入参
func GetFlopsInfo(gpuName string) InfoFlop {
	logs.Normal("Getting FLOPS info...")

	gpuInfos := gpu.InitGpuInfos()
	for _, info := range gpuInfos {
        if strings.Contains(info.Name, gpuName) {
            return InfoFlop{Flops: info.Fp32}
        }
    }
	return InfoFlop{Flops: 0}
}