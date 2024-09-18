package main

import (
	"fmt"
	"log"
	"mymanagerdocker/application/api"
	"mymanagerdocker/application/service"
	"mymanagerdocker/application/web"
	"net/http"
)

func main() {

	// services
	docker, _ := service.Build()

	// api
	api.Build(docker)

	// web
	web.Build(docker)

	// Porta 8080 para o servidor HTTP
	port := ":8080"
	fmt.Printf("Servidor rodando na porta %s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
