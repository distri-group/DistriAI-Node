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
		GpuInfo{Name: "NVIDIA GeForce RTX 4070 Ti SUPER", Fp16: 44.10, Fp32: 44.10, Fp64: 0.689},
		GpuInfo{Name: "NVIDIA GeForce GTX 1080", Fp16: 0.138, Fp32: 8.87, Fp64: 0.277},
		GpuInfo{Name: "NVIDIA GeForce GTX 1080 11Gbps", Fp16: 0.138, Fp32: 8.87, Fp64: 0.277},
		GpuInfo{Name: "NVIDIA GeForce GTX 1080 Max-Q", Fp16: 0.109, Fp32: 6.99, Fp64: 0.218},
		GpuInfo{Name: "NVIDIA GeForce GTX 1080 Mobile", Fp16: 0.138, Fp32: 8.87, Fp64: 0.277},
		GpuInfo{Name: "NVIDIA GeForce GTX 1080 Ti", Fp16: 0.177, Fp32: 11.34, Fp64: 0.354},
		GpuInfo{Name: "NVIDIA GeForce GTX 1080 Ti 10 GB", Fp16: 0.167, Fp32: 10.69, Fp64: 0.334},
		GpuInfo{Name: "NVIDIA GeForce GTX 1080 Ti 12 GB", Fp16: 0.167, Fp32: 10.69, Fp64: 0.334},
		GpuInfo{Name: "NVIDIA Tesla C1080", Fp16: 0.0, Fp32: 0.622, Fp64: 0.077})
	return gpuInfos
}
