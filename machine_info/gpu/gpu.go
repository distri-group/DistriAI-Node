package gpu

import (
	"DistriAI-Node/pattern"
	"DistriAI-Node/utils/log_utils"
	"bytes"
	"os/exec"
	"strconv"
	"strings"
)

// InfoGPU 定义 InfoGPU 结构体
type InfoGPU struct {
	Model  string `json:"Model"`  // GPU显卡型号
	Number int    `json:"Number"` // GPU显卡数量
}

// GetIntelGPUInfo 获取 Intel GPU 信息并返回一个包含 InfoGPU 结构体的切片
func GetIntelGPUInfo() (InfoGPU, error) {
	logs.Normal("Getting GPU info...")

	var gpuInfo InfoGPU

	cmd := exec.Command("nvidia-smi", "--query-gpu=count,gpu_name", "--format=csv,noheader")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		// Easy debugging
		return InfoGPU{Model: pattern.NO_GPU, Number: 0}, err
		// return InfoGPU{Model: "NVIDIA RTX 3090", Number: 20}, err
	}

	result := strings.Split(strings.TrimSpace(out.String()), ", ")

	if len(result) >= 2 {

		number, err := strconv.Atoi(result[0])
		if err != nil {
			number = 0
		}

		gpuInfo = InfoGPU{
			Model:  result[1],
			Number: number,
		}
	}
	return gpuInfo, nil
}
