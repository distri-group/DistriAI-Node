package disk

import (
	"DistriAI-Node/utils/log_utils"
	"fmt"

	"github.com/shirou/gopsutil/v3/disk"
)

// InfoDisk 定义 InfoDisk 结构体
type InfoDisk struct {
	Path       string  `json:"Path"`       // 磁盘路径（如：/dev/sda1）
	TotalSpace float64 `json:"TotalSpace"` // 总共的硬盘空间大小，单位 GB
	// FreeSpace  float64 `json:"FreeSpace"`  // 可用的硬盘空间大小，单位 GB
}

// GetDiskInfo 获取硬盘信息并返回 InfoDisk 结构体切片
func GetDiskInfo() (InfoDisk, error) {
	partitions, err := disk.Partitions(true)
	if err != nil {
		return InfoDisk{}, fmt.Errorf("failed to get disk info: %w", err)
	}

	var max uint64 = 0
	var largestPartition disk.PartitionStat

	for _, partition := range partitions {
		usageStat, err := disk.Usage(partition.Mountpoint)
		if err != nil {
			logs.Error(fmt.Sprintf("get %v usage stat failed, err:%v\n", partition.Mountpoint, err))
			continue
		}

		if usageStat.Total > max {
			max = usageStat.Total
			largestPartition = partition
		}
	}

	diskInfo := InfoDisk{
		Path:       largestPartition.Device,
		TotalSpace: float64(max) / 1024 / 1024 / 1024,
	}

	// diskInfos := make([]InfoDisk, 0)
	// for _, p := range partitions {
	// 	usage, err := disk.Usage(p.Mountpoint)
	// 	if err != nil {
	// 		continue
	// 	}
	// 	diskInfos = append(diskInfos, InfoDisk{
	// 		Path:       p.Device,
	// 		TotalSpace: float64(usage.Total) / (1 << 30),
	// 		FreeSpace:  float64(usage.Free) / (1 << 30),
	// 	})
	// }
	return diskInfo, nil
}
