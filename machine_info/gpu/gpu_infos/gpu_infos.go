package gpu

type GpuInfo struct {
	Name string
	Fp16 float32
	Fp32 float32
	Fp64 float32
}

func InitGpuInfos() []GpuInfo {
	var gpuInfos []GpuInfo
	gpuInfos = append(gpuInfos,
		GpuInfo{Name: "NVIDIA GeForce RTX 3090", Fp16: 35.58, Fp32: 35.58, Fp64: 0.556},
		GpuInfo{Name: "NVIDIA GeForce RTX 3090 Ti", Fp16: 40.00, Fp32: 40.00, Fp64: 0.625},
		GpuInfo{Name: "NVIDIA GeForce RTX 3050 6 GB", Fp16: 6.021, Fp32: 6.021, Fp64: 0.094},
		GpuInfo{Name: "NVIDIA GeForce RTX 4070 SUPER", Fp16: 35.48, Fp32: 35.48, Fp64: 0.554},
		GpuInfo{Name: "NVIDIA GeForce RTX 4070 Ti SUPER", Fp16: 44.10, Fp32: 44.10, Fp64: 0.689})
	return gpuInfos
}
