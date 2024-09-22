package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"mymanagerdocker/application/model"
	"mymanagerdocker/application/util"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
)

type DockerService struct {
	Cli     *client.Client
	Ctx     context.Context
	Host    string
	Version string
}

func Build() (DockerService, error) {
	Cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return DockerService{}, fmt.Errorf("erro ao criar cliente Docker: %v", err)
	}

	Ctx := context.Background()

	Host := os.Getenv("DOCKER_HOST") // Host de conexão (ex: unix:///var/run/docker.sock)
	if Host == "" {
		Host = "unix:///var/run/docker.sock" // Valor padrão se DOCKER_HOST não estiver definido
	}

	Version := Cli.ClientVersion() // Versão da API Docker que o cliente está usando

	// Exibir as configurações
	fmt.Printf("Conectando ao Docker no host: %s\n", Host)
	fmt.Printf("Versão da API Docker: %s\n", Version)

	return DockerService{
		Cli,
		Ctx,
		Host,
		Version,
	}, nil
}

func (s *DockerService) GetContainers() ([]types.Container, error) {

	containers, err := s.Cli.ContainerList(s.Ctx, container.ListOptions{})

	if err != nil {
		log.Fatalf("Erro ao listar containers: %v", err)
	}

	return containers, err
}

func (s *DockerService) GetNetworks() ([]model.NetworkInfo, error) {
	networks, err := s.Cli.NetworkList(s.Ctx, network.ListOptions{})

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
	containers, err := s.Cli.ContainerList(s.Ctx, container.ListOptions{})
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
	stats, err := s.Cli.ContainerStats(s.Ctx, containerID, false)

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

// Função para iniciar um container
func (s *DockerService) StartContainer(containerID string) error {

	if err := s.Cli.ContainerStart(s.Ctx, containerID, container.StartOptions{}); err != nil {
		return fmt.Errorf("failed to start container: %w", err)
	}
	return nil
}

// Função para parar um container
func (s *DockerService) StopContainer(containerID string) error {

	if err := s.Cli.ContainerStop(s.Ctx, containerID, container.StopOptions{}); err != nil {
		return fmt.Errorf("failed to stop container: %w", err)
	}
	return nil
}

// Função para duplicar um container
func (s *DockerService) DuplicateContainer(containerID string) (string, error) {

	// Obter as informações do container original
	containerJSON, err := s.Cli.ContainerInspect(s.Ctx, containerID)
	if err != nil {
		return "", fmt.Errorf("failed to inspect container: %w", err)
	}

	// Criar uma nova configuração baseada na original
	newContainerConfig := containerJSON.Config
	newContainer, err := s.Cli.ContainerCreate(s.Ctx, newContainerConfig, nil, nil, nil, "")
	if err != nil {
		return "", fmt.Errorf("failed to create duplicate container: %w", err)
	}

	return newContainer.ID, nil
}
