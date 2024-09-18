package response

import (
	"mymanagerdocker/application/model"
)

type DockerInfo struct {
	ContainerStats []model.ContainerInfo
	Networks       []model.NetworkInfo
}
