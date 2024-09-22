package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func main2() {

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Fatalf("erro ao criar cliente Docker: %v", err)
	}

	ctx := context.Background()

	containerID := "51e892815603b12271040cacfd2dbe9c3d28e0647a2ef1bb09245bacd4a790d7"

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

	// Goroutine para capturar entrada do usuário (os.Stdin)
	go func() {
		for {

			// Buffer temporário para capturar a entrada do usuário
			buf := make([]byte, 1024)
			n, err := os.Stdin.Read(buf)
			if err != nil && err != io.EOF {
				log.Fatal(err)
			}

			// Envia a entrada do usuário para o contêiner
			// Indicador de que estamos esperando por entrada do usuário
			fmt.Print("<<SEND ")
			os.Stdout.Write(buf[:n])

			_, err = execStartResp.Conn.Write(buf[:n])

			if err != nil {
				log.Fatal(err)
			}

		}
	}()

	// Goroutine para gerenciar a saída do contêiner (stdout e stderr)
	go func() {
		// Buffer para capturar a saída
		bufnew := make([]byte, 1024)
		for {
			// Lê a saída do contêiner
			n, err := execStartResp.Reader.Read(bufnew)
			if err != nil && err != io.EOF {
				log.Fatal(err)
			}
			if n > 0 {
				// Indicador de que recebemos saída do contêiner
				fmt.Print(">>RECIVE ")

				// Escreve a saída no stdout
				os.Stdout.Write(bufnew[:n])
			}
		}
	}()

	// Mantém o processo vivo
	select {}
}
