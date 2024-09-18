package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"mymanagerdocker/application/model"
	"mymanagerdocker/application/util"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
)

type DockerService struct {
	cli *client.Client
	ctx context.Context
}

func Build() (DockerService, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return DockerService{}, fmt.Errorf("erro ao criar cliente Docker: %v", err)
	}

	ctx := context.Background()

	return DockerService{
		cli,
		ctx,
	}, nil
}

func (s *DockerService) GetContainers() ([]types.Container, error) {

	containers, err := s.cli.ContainerList(s.ctx, container.ListOptions{})

	if err != nil {
		log.Fatalf("Erro ao listar containers: %v", err)
	}

	return containers, err
}

func (s *DockerService) GetNetworks() ([]model.NetworkInfo, error) {
	networks, err := s.cli.NetworkList(s.ctx, network.ListOptions{})

	if err != nil {
		log.Fatalf("Erro ao listar redes: %v", err)
	}

	networkStats := make([]model.NetworkInfo, 0)

	for _, network := range networks {
		networkStats = append(networkStats,
			model.NetworkInfo{
				Network: network,
			})
	}

	return networkStats, err
}

func (s *DockerService) GetContainersStats() ([]model.ContainerInfo, error) {
	containers, err := s.cli.ContainerList(s.ctx, container.ListOptions{})
	if err != nil {
		log.Fatalf("Erro ao listar containers: %v", err)
		return nil, err
	}

	containersStats := make([]model.ContainerInfo, 0)

	for _, container := range containers {
		stats, err := s.GetContainerStats(container.ID)
		if err != nil {
			return nil, fmt.Errorf("erro ao obter estatísticas do container %s: %v", container.ID, err)
		}
		memUsage, memLimit, memPercent := util.CalculateMemoryPercent(stats)

		containersStats = append(containersStats, model.ContainerInfo{
			Container:        container,
			Stats:            stats,
			CPUPercentage:    util.CalculateCPUUsage(stats),
			MemoryUsage:      memUsage,
			MemoryLimit:      memLimit,
			MemoryPercentage: memPercent,
		})
	}

	return containersStats, nil
}

func (s *DockerService) GetContainerStats(containerID string) (container.StatsResponse, error) {
	stats, err := s.cli.ContainerStats(s.ctx, containerID, false)

	if err != nil {
		return container.StatsResponse{}, fmt.Errorf("erro ao obter estatísticas do container %s: %v", containerID, err)
	}

	defer stats.Body.Close()

	data, err := ioutil.ReadAll(stats.Body)
	if err != nil {
		return container.StatsResponse{}, err
	}

	// Decodifica os dados de estatísticas em JSON
	var statsJSON container.StatsResponse
	if err := json.Unmarshal(data, &statsJSON); err != nil {
		return container.StatsResponse{}, err
	}

	return statsJSON, nil
}
