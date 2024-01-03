package flops

import (
	"DistriAI-Node/machine_info/gpu/gpu_infos"
	"DistriAI-Node/utils/log_utils"
	"strings"
)

type InfoFlop struct {
	Flops float32 `json:"Flops"`
}

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
