package api

import (
	"fmt"
	"mymanagerdocker/application/api/response"
	"mymanagerdocker/application/service"
	"mymanagerdocker/application/util"
	"net/http"
)

type ApiData struct {
	service service.DockerService
}

func Build(s service.DockerService) ApiData {
	apiData := ApiData{
		service: s,
	}

	http.HandleFunc("/api", apiData.serve)
	http.HandleFunc("/api/containers", apiData.containers)
	http.HandleFunc("/api/networks", apiData.networks)

	return apiData
}

func (a *ApiData) serve(w http.ResponseWriter, r *http.Request) {
	containerStats, e := a.service.GetContainersStats()

	if e != nil {
		http.Error(w, fmt.Sprintf("Erro ao codificar JSON: %v", e), http.StatusInternalServerError)
		return
	}

	networks, e := a.service.GetNetworks()

	util.SendJson(w, response.DockerInfo{
		ContainerStats: containerStats,
		Networks:       networks,
	})
}

func (a *ApiData) containers(w http.ResponseWriter, r *http.Request) {
	containers, e := a.service.GetContainers()

	if e != nil {
		http.Error(w, fmt.Sprintf("Erro ao codificar JSON: %v", e), http.StatusInternalServerError)
		return
	}

	util.SendJson(w, containers)
}

func (a *ApiData) networks(w http.ResponseWriter, r *http.Request) {
	networks, e := a.service.GetNetworks()

	if e != nil {
		http.Error(w, fmt.Sprintf("Erro ao codificar JSON: %v", e), http.StatusInternalServerError)
		return
	}

	util.SendJson(w, networks)
}
