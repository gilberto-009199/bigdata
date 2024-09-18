package model

import (
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
)

type ContainerInfo struct {
	Container        types.Container         `json:"container"`
	Stats            container.StatsResponse `json:"stats"`
	CPUPercentage    float64                 `json:"cpu_percentage"`
	MemoryUsage      uint64                  `json:"memory_usage"`
	MemoryLimit      uint64                  `json:"memory_limit"`
	MemoryPercentage float64                 `json:"memory_percentage"`
}
