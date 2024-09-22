package api

import (
	"archive/tar"
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"mymanagerdocker/application/api/response"
	"mymanagerdocker/application/service"
	"mymanagerdocker/application/util"
	"net/http"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	// CORS PERMIT ALL *.*
	CheckOrigin: func(r *http.Request) bool { return true },
}

type ApiData struct{ service service.DockerService }

func Build(s service.DockerService) ApiData {
	apiData := ApiData{service: s}

	http.HandleFunc("/api", apiData.serve)
	http.HandleFunc("/api/containers", apiData.containers)
	http.HandleFunc("/api/networks", apiData.networks)
	http.HandleFunc("/upload", apiData.upload)
	http.HandleFunc("/ws", apiData.sockTerminal)

	http.HandleFunc("/api/start/", apiData.startHandler)
	http.HandleFunc("/api/stop/", apiData.stopHandler)
	http.HandleFunc("/api/duplicar/", apiData.duplicateHandler)

	return apiData
}

func (a *ApiData) startHandler(w http.ResponseWriter, r *http.Request) {
	containerId := r.URL.Query().Get("containerId")
	a.service.StartContainer(containerId)
}

func (a *ApiData) stopHandler(w http.ResponseWriter, r *http.Request) {
	containerId := r.URL.Query().Get("containerId")
	a.service.StopContainer(containerId)
}

func (a *ApiData) duplicateHandler(w http.ResponseWriter, r *http.Request) {
	containerId := r.URL.Query().Get("containerId")
	a.service.DuplicateContainer(containerId)
}

// Web Socket
func (a *ApiData) sockTerminal(w http.ResponseWriter, r *http.Request) {
	// Atualiza a conexão HTTP para WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Erro ao estabelecer conexão WebSocket:", err)
		return
	}
	defer conn.Close()
	// code de send and recive std stout container
	// /ws?containerId=<????>
	containerID := r.URL.Query().Get("containerId")
	if containerID == "" {
		log.Println("containerId não fornecido")
		conn.WriteMessage(websocket.TextMessage, []byte("containerId não fornecido"))
		return
	}
	// @todo pensar em colocar um tipo de callback , para manter a logica no DockerService
	cli := a.service.Cli
	ctx := a.service.Ctx

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

func (a *ApiData) serve(w http.ResponseWriter, r *http.Request) {
	containerStats, e := a.service.GetContainersStats()
	if e != nil {
		http.Error(w, fmt.Sprintf("Erro ao codificar JSON: %v", e), http.StatusInternalServerError)
		return
	}

	networks, e := a.service.GetNetworks()
	if e != nil {
		http.Error(w, fmt.Sprintf("Erro ao pegar as Networks : %v", e), http.StatusInternalServerError)
		return
	}

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

func (a *ApiData) upload(w http.ResponseWriter, r *http.Request) {

	// Limita o tamanho do arquivo que pode ser enviado (exemplo: 10 MB)
	r.ParseMultipartForm(10 << 20)

	// Obtém o arquivo e as informações
	file, handler, err := r.FormFile("file")
	if err != nil {
		fmt.Println("Erro ao receber o arquivo:", err)
		http.Error(w, "Erro ao receber o arquivo", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Obtém o ID do container
	containerId := r.FormValue("containerId")
	if containerId == "" {
		http.Error(w, "ID do container não fornecido", http.StatusBadRequest)
		return
	}

	// Salvando o arquivo em uma pasta temporária
	dst, err := os.Create(fmt.Sprintf("/tmp/%s", handler.Filename))
	if err != nil {
		http.Error(w, "Erro ao salvar o arquivo", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	// Copia o conteúdo do arquivo
	_, err = io.Copy(dst, file)
	if err != nil {
		http.Error(w, "Erro ao salvar o arquivo", http.StatusInternalServerError)
		return
	}

	// Agora, enviaremos o arquivo para o container Docker
	// @todo pensar em colocar um tipo de callback , para manter a logica no DockerService
	cli := a.service.Cli
	ctx := a.service.Ctx

	err = copyFileToContainer(cli, ctx, containerId, dst.Name(), "/")
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Erro ao enviar o arquivo para o container", http.StatusInternalServerError)
		return
	}

	// Envia a resposta de sucesso
	w.Write([]byte("Upload feito com sucesso"))
}

// Função que cria um arquivo tar para enviar ao container
func createTarFile(sourceFilePath string) (*bytes.Buffer, error) {
	buf := new(bytes.Buffer)
	tw := tar.NewWriter(buf)

	// Abre o arquivo de origem
	file, err := os.Open(sourceFilePath)
	if err != nil {
		return nil, fmt.Errorf("erro ao abrir arquivo de origem: %v", err)
	}
	defer file.Close()

	// Obtém informações do arquivo
	info, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("erro ao obter informações do arquivo: %v", err)
	}

	// Escreve o header TAR
	header := &tar.Header{
		Name: info.Name(),
		Size: info.Size(),
		Mode: 0600,
	}

	err = tw.WriteHeader(header)
	if err != nil {
		return nil, fmt.Errorf("erro ao escrever cabeçalho TAR: %v", err)
	}

	// Copia o conteúdo do arquivo para o TAR
	_, err = io.Copy(tw, file)
	if err != nil {
		return nil, fmt.Errorf("erro ao copiar arquivo para TAR: %v", err)
	}

	// Fecha o writer do TAR
	err = tw.Close()
	if err != nil {
		return nil, fmt.Errorf("erro ao fechar TAR: %v", err)
	}

	return buf, nil
}

// Função que envia o arquivo tar para o container
func copyFileToContainer(cli *client.Client, ctx context.Context, containerID, srcPath, destPath string) error {

	// Cria um arquivo tar com o conteúdo a ser enviado
	tarReader, err := createTarFile(srcPath)
	if err != nil {
		return fmt.Errorf("erro ao criar arquivo tar: %v", err)
	}
	fmt.Println("CopyToContainer!!!!")
	// Copia o arquivo tar para o container
	err = cli.CopyToContainer(ctx, containerID, destPath, tarReader, container.CopyToContainerOptions{})
	if err != nil {
		return fmt.Errorf("erro ao copiar arquivo para o container: %v", err)
	}

	fmt.Printf("Arquivo %s enviado para o container %s no diretório %s\n", srcPath, containerID, destPath)
	return nil
}
