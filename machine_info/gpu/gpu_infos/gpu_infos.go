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
		GpuInfo{Name: "NVIDIA Tesla C1080", Fp16: 0.0, Fp32: 0.622, Fp64: 0.077},
		GpuInfo{Name: "NVIDIA GeForce RTX 4080 SUPER", Fp16: 51.3, Fp32: 51.3, Fp64: 0.801},
		GpuInfo{Name: "NVIDIA GeForce RTX 4090D", Fp16: 73.54, Fp32: 73.54, Fp64: 1.149},
		GpuInfo{Name: "NVIDIA GeForce RTX 4060 AD106", Fp16: 19.47, Fp32: 19.47, Fp64: 0.304},
		GpuInfo{Name: "NVIDIA GeForce RTX 4080 Ti", Fp16: 67.58, Fp32: 67.58, Fp64: 1.056},
		GpuInfo{Name: "NVIDIA GeForce RTX 4050", Fp16: 13.52, Fp32: 13.52, Fp64: 0.211},
		GpuInfo{Name: "NVIDIA RTX 4000 Ada Generation", Fp16: 26.73, Fp32: 26.73, Fp64: 0.417},
		GpuInfo{Name: "NVIDIA RTX 4500 Ada Generation", Fp16: 39.63, Fp32: 39.63, Fp64: 0.619},
		GpuInfo{Name: "NVIDIA RTX 5000 Ada Generation", Fp16: 65.28, Fp32: 65.28, Fp64: 1.02},
		GpuInfo{Name: "NVIDIA GeForce RTX 4060", Fp16: 15.11, Fp32: 15.11, Fp64: 0.236},
		GpuInfo{Name: "NVIDIA GeForce RTX 4060 Ti 16 GB", Fp16: 22.06, Fp32: 22.06, Fp64: 0.344},
		GpuInfo{Name: "NVIDIA GeForce RTX 4060 Ti 8 GB", Fp16: 22.06, Fp32: 22.06, Fp64: 0.344},
		GpuInfo{Name: "NVIDIA GeForce RTX 4070", Fp16: 29.15, Fp32: 29.15, Fp64: 0.455},
		GpuInfo{Name: "NVIDIA H100 CNX", Fp16: 215.4, Fp32: 53.84, Fp64: 26.92},
		GpuInfo{Name: "NVIDIA H100 PCIe 80 GB", Fp16: 204.9, Fp32: 51.22, Fp64: 25.61},
		GpuInfo{Name: "NVIDIA H100 PCIe 96 GB", Fp16: 248.3, Fp32: 62.08, Fp64: 31.04},
		GpuInfo{Name: "NVIDIA H100 SXM5 64 GB", Fp16: 267.6, Fp32: 66.91, Fp64: 33.45},
		GpuInfo{Name: "NVIDIA H100 SXM5 80 GB", Fp16: 267.6, Fp32: 66.91, Fp64: 33.45},
		GpuInfo{Name: "NVIDIA H100 SXM5 96 GB", Fp16: 248.3, Fp32: 62.08, Fp64: 31.04},
		GpuInfo{Name: "NVIDIA H800 PCIe 80 GB", Fp16: 204.9, Fp32: 51.22, Fp64: 25.61},
		GpuInfo{Name: "NVIDIA H800 SXM5", Fp16: 237.2, Fp32: 59.3, Fp64: 29.65},
		GpuInfo{Name: "NVIDIA L4", Fp16: 31.33, Fp32: 31.33, Fp64: 0.489},
		GpuInfo{Name: "NVIDIA RTX 2000 Embedded Ada Generation", Fp16: 12.99, Fp32: 12.99, Fp64: 0.203},
		GpuInfo{Name: "NVIDIA RTX 2000 Max-Q Ada Generation", Fp16: 8.94, Fp32: 8.94, Fp64: 0.139},
		GpuInfo{Name: "NVIDIA RTX 2000 Mobile Ada Generation", Fp16: 12.99, Fp32: 12.99, Fp64: 0.203},
		GpuInfo{Name: "NVIDIA RTX 3000 Mobile Ada Generation", Fp16: 15.62, Fp32: 15.62, Fp64: 0.244},
		GpuInfo{Name: "NVIDIA RTX 3500 Mobile Ada Generation", Fp16: 15.82, Fp32: 15.82, Fp64: 0.247},
		GpuInfo{Name: "NVIDIA RTX 4000 Mobile Ada Generation", Fp16: 24.72, Fp32: 24.72, Fp64: 0.386},
		GpuInfo{Name: "NVIDIA RTX 4000 SFF Ada Generation", Fp16: 19.17, Fp32: 19.17, Fp64: 0.299},
		GpuInfo{Name: "NVIDIA RTX 5000 Max-Q Ada Generation", Fp16: 32.69, Fp32: 32.69, Fp64: 0.51},
		GpuInfo{Name: "NVIDIA RTX 5000 Mobile Ada Embedded", Fp16: 41.15, Fp32: 41.15, Fp64: 0.643},
		GpuInfo{Name: "NVIDIA RTX 5000 Mobile Ada Generation", Fp16: 41.15, Fp32: 41.15, Fp64: 0.643},
		GpuInfo{Name: "NVIDIA Jetson AGX Orin 64 GB", Fp16: 10.65, Fp32: 5.32, Fp64: 0.002},
		GpuInfo{Name: "NVIDIA Jetson Orin NX 8 GB", Fp16: 3.13, Fp32: 1.56, Fp64: 0.783},
		GpuInfo{Name: "NVIDIA Jetson Orin Nano 4 GB", Fp16: 1.28, Fp32: 0.64, Fp64: 0.32},
		GpuInfo{Name: "NVIDIA Jetson Orin Nano 8 GB", Fp16: 2.56, Fp32: 1.28, Fp64: 0.64},
		GpuInfo{Name: "NVIDIA Jetson AGX Orin 32 GB", Fp16: 6.66, Fp32: 3.33, Fp64: 1.66},
		GpuInfo{Name: "NVIDIA Jetson Orin NX 16 GB", Fp16: 3.76, Fp32: 1.88, Fp64: 0.94},
		GpuInfo{Name: "NVIDIA GeForce RTX 4090 Ti", Fp16: 93.24, Fp32: 93.24, Fp64: 1.457},
		GpuInfo{Name: "NVIDIA TITAN Ada", Fp16: 92.9, Fp32: 92.9, Fp64: 1.452})
	return gpuInfos
}
