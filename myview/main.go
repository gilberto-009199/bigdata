package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
)

type DockerInfo struct {
	Containers []ContainerInfo                  `json:"containers"`
	Networks   map[string]types.NetworkResource `json:"networks"`
}

type ContainerInfo struct {
	Container types.Container `json:"container"`
	Stats     ContainerStats  `json:"stats"`
}

type ContainerStats struct {
	CPUPercentage    float64                           `json:"cpu_percentage"`
	MemoryUsage      uint64                            `json:"memory_usage"`
	MemoryLimit      uint64                            `json:"memory_limit"`
	MemoryPercentage float64                           `json:"memory_percentage"`
	Networks         map[string]container.NetworkStats `json:"networks"`
}

func calculateCPUUsage(stats *types.StatsJSON) float64 {

	cpuDelta := float64(stats.CPUStats.CPUUsage.TotalUsage) - float64(stats.PreCPUStats.CPUUsage.TotalUsage)
	systemDelta := float64(stats.CPUStats.SystemUsage) - float64(stats.PreCPUStats.SystemUsage)
	numberOfCores := float64(stats.CPUStats.OnlineCPUs)

	return (cpuDelta / systemDelta) * numberOfCores * 100.0
}

func calculateMemoryPercent(stats *types.StatsJSON) (usage uint64, limit uint64, percentage float64) {
	usage = stats.MemoryStats.Usage
	limit = stats.MemoryStats.Limit
	if limit > 0 {
		percentage = (float64(usage) / float64(limit)) * 100.0
	}
	return
}

// Função para obter estatísticas de um container
func getContainerStats(cli *client.Client, containerID string) (ContainerStats, error) {
	ctx := context.Background()
	stats, err := cli.ContainerStats(ctx, containerID, false)
	if err != nil {
		return ContainerStats{}, err
	}
	defer stats.Body.Close()

	// Lê os dados de estatísticas
	data, err := ioutil.ReadAll(stats.Body)
	if err != nil {
		return ContainerStats{}, err
	}

	// Decodifica os dados de estatísticas em JSON
	var statsJSON types.StatsJSON
	if err := json.Unmarshal(data, &statsJSON); err != nil {
		return ContainerStats{}, err
	}

	// Calcula o uso de CPU e Memória
	cpuUsage := calculateCPUUsage(&statsJSON)
	memUsage, memLimit, memPercent := calculateMemoryPercent(&statsJSON)

	// Calcule o uso da CPU (percentual instantâneo aproximado)
	numCPU := float64(len(statsJSON.CPUStats.CPUUsage.PercpuUsage))
	if numCPU == 0 {
		numCPU = 1 // Se não tiver dados sobre CPUs, assume que há pelo menos uma CPU
	}

	return ContainerStats{
		CPUPercentage:    cpuUsage,
		MemoryUsage:      memUsage,
		MemoryLimit:      memLimit,
		MemoryPercentage: memPercent,
		Networks:         statsJSON.Networks,
	}, nil
}

// Função para obter informações de containers e networks
func getDockerInfo() (DockerInfo, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return DockerInfo{}, fmt.Errorf("erro ao criar cliente Docker: %v", err)
	}

	ctx := context.Background()

	// Listar containers
	containers, err := cli.ContainerList(ctx, container.ListOptions{})
	if err != nil {
		return DockerInfo{}, fmt.Errorf("erro ao listar containers: %v", err)
	}

	// Listar networks
	networks, err := cli.NetworkList(ctx, network.ListOptions{})
	if err != nil {
		return DockerInfo{}, fmt.Errorf("erro ao listar networks: %v", err)
	}

	// Organizar as informações em uma estrutura
	info := DockerInfo{
		Containers: []ContainerInfo{},
		Networks:   make(map[string]network.Inspect),
	}

	// Adicionar estatísticas para cada container
	for _, container := range containers {
		stats, err := getContainerStats(cli, container.ID)
		if err != nil {
			log.Printf("Erro ao obter estatísticas para o container %s: %v", container.ID, err)
		}

		info.Containers = append(info.Containers, ContainerInfo{
			Container: container,
			Stats:     stats,
		})
	}

	for _, network := range networks {
		info.Networks[network.ID] = network
	}

	return info, nil
}

// Handler para a API que retorna informações sobre containers e networks
func dockerInfoAPI(w http.ResponseWriter, r *http.Request) {
	info, err := getDockerInfo()
	if err != nil {
		http.Error(w, fmt.Sprintf("Erro ao obter informações do Docker: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(info); err != nil {
		http.Error(w, fmt.Sprintf("Erro ao codificar JSON: %v", err), http.StatusInternalServerError)
	}
}

// Handler para servir o arquivo index.html
func serveIndex(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/index.html")
}

func main() {
	// Roteamento
	http.HandleFunc("/", serveIndex)                   // Página inicial
	http.HandleFunc("/api/docker-info", dockerInfoAPI) // API

	// Porta 8080 para o servidor HTTP
	port := ":8080"
	fmt.Printf("Servidor rodando na porta %s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
