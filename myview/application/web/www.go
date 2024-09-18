package web

import (
	"mymanagerdocker/application/service"
	"net/http"
	"path/filepath"
	"strings"
)

type WebData struct {
	service service.DockerService
}

func Build(s service.DockerService) WebData {

	webData := WebData{
		service: s,
	}

	// Serve arquivos estáticos
	http.HandleFunc("/", webData.serve)

	return webData
}

func (a *WebData) serve(w http.ResponseWriter, r *http.Request) {
	// Se a rota for "/" ou "/index.html", serve o arquivo HTML
	if r.URL.Path == "/" {
		http.ServeFile(w, r, "static/index.html")
		return
	}

	// Serve arquivos estáticos para todas as outras rotas
	staticFile := filepath.Join("static", r.URL.Path)

	// Verifica se o caminho solicitado é um arquivo .js, .css, etc.
	if strings.HasSuffix(r.URL.Path, ".js") || strings.HasSuffix(r.URL.Path, ".css") || strings.HasSuffix(r.URL.Path, ".png") {
		http.ServeFile(w, r, staticFile)
	} else {
		http.NotFound(w, r)
	}
}
