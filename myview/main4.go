package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// Permite conexões de qualquer origem (apenas para desenvolvimento)
		return true
	},
}

func handleTerminal(w http.ResponseWriter, r *http.Request) {
	// Atualiza a conexão HTTP para WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Erro ao estabelecer conexão WebSocket:", err)
		return
	}
	defer conn.Close()

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Fatalf("erro ao criar cliente Docker: %v", err)
	}

	ctx := context.Background()
	containerID := "51e892815603b12271040cacfd2dbe9c3d28e0647a2ef1bb09245bacd4a790d7"

	// Cria a configuração do comando `sh` no contêiner
	execConfig := types.ExecConfig{
		Cmd:          []string{"sh"},
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
		Tty:          true,
	}

	execIDResp, err := cli.ContainerExecCreate(ctx, containerID, execConfig)
	if err != nil {
		log.Fatal(err)
	}

	execStartCheck := types.ExecStartCheck{}
	execStartResp, err := cli.ContainerExecAttach(ctx, execIDResp.ID, execStartCheck)
	if err != nil {
		log.Fatal(err)
	}
	defer execStartResp.Close()

	// Goroutine para capturar entrada do cliente via WebSocket e enviar ao contêiner
	go func() {
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				log.Println("Erro ao ler mensagem WebSocket:", err)
				break
			}

			// Indicador de que estamos enviando a entrada ao contêiner
			fmt.Print("<<SEND ")

			// Envia o comando do WebSocket ao contêiner Docker
			_, err = execStartResp.Conn.Write(message)
			if err != nil {
				log.Println("Erro ao enviar comando para o contêiner:", err)
				break
			}
		}
	}()

	// Loop para capturar a saída do contêiner e enviar de volta ao cliente WebSocket
	buf := make([]byte, 1024)
	for {
		n, err := execStartResp.Reader.Read(buf)
		if err != nil && err != io.EOF {
			log.Println("Erro ao ler saída do contêiner:", err)
			break
		}
		if n > 0 {
			// Indicador de que recebemos saída do contêiner
			fmt.Print(">>RECIVE ")

			// Envia a saída do contêiner de volta ao cliente WebSocket como mensagem binária
			err = conn.WriteMessage(websocket.BinaryMessage, buf[:n])
			if err != nil {
				log.Println("Erro ao enviar saída para o WebSocket:", err)
				break
			}
		}
	}
}

func main4() {
	http.HandleFunc("/ws", handleTerminal)
	log.Println("Servidor WebSocket iniciado em :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
