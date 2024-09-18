package model

import (
	"github.com/docker/docker/api/types/network"
)

type NetworkInfo struct {
	Network network.Inspect `json:"network"`
}
