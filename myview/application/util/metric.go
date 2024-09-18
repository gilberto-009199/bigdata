package util

import "github.com/docker/docker/api/types/container"

func CalculateCPUUsage(stats container.StatsResponse) float64 {

	cpuDelta := float64(stats.CPUStats.CPUUsage.TotalUsage) - float64(stats.PreCPUStats.CPUUsage.TotalUsage)
	systemDelta := float64(stats.CPUStats.SystemUsage) - float64(stats.PreCPUStats.SystemUsage)
	numberOfCores := float64(stats.CPUStats.OnlineCPUs)

	return (cpuDelta / systemDelta) * numberOfCores * 100.0
}

func CalculateMemoryPercent(stats container.StatsResponse) (usage uint64, limit uint64, percentage float64) {
	usage = stats.MemoryStats.Usage
	limit = stats.MemoryStats.Limit
	if limit > 0 {
		percentage = (float64(usage) / float64(limit)) * 100.0
	}
	return
}
