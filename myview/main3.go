package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/gorilla/websocket"
)

/*var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}*/

func handleConnection(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Failed to establish WebSocket connection:", err)
		return
	}
	defer conn.Close()

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())

	if err != nil {
		fmt.Errorf("erro ao criar cliente Docker: %v", err)
	}

	ctx := context.Background()

	var containerID string
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
			break
		}

		if containerID == "" {
			containerID = string(msg)
			continue
		}

		execConfig := types.ExecConfig{
			Cmd:          []string{"sh"},
			AttachStdin:  true,
			AttachStdout: true,
			AttachStderr: true,
			Tty:          true,
		}

		execIDResp, err := cli.ContainerExecCreate(ctx, containerID, execConfig)
		if err != nil {
			log.Println("Exec create error:", err)
			continue
		}

		execStartCheck := types.ExecStartCheck{}
		execStartResp, err := cli.ContainerExecAttach(ctx, execIDResp.ID, execStartCheck)
		if err != nil {
			log.Println("Exec start error:", err)
			continue
		}

		// Use bufio.Reader to handle output
		reader := bufio.NewReader(execStartResp.Reader)
		go func() {
			defer execStartResp.Close()

			for {
				line, err := reader.Peek(1024)
				if err != nil {
					if err != io.EOF {
						log.Println("Error reading output:", err)
					}
					break
				}

				err = conn.WriteMessage(websocket.TextMessage, line)
				if err != nil {
					log.Println("Error writing message to WebSocket:", err)
					break
				}
			}
		}()
	}
}

func main3() {
	http.HandleFunc("/ws", handleConnection)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
