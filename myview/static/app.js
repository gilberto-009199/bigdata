var app = {
    upload(containerId) {
        const fileInput = document.getElementById('fileInput');

        // Simula um clique para abrir o seletor de arquivos
        fileInput.click();

        // Quando o usuário selecionar um arquivo
        fileInput.onchange = function(event) {
            const file = event.target.files[0];

            if (!file) {
                alert("Nenhum arquivo selecionado!");
                return;
            }

            // Cria um formData para enviar o arquivo
            const formData = new FormData();
            formData.append('file', file);
            formData.append('containerId', containerId);

            // Faz o upload via AJAX (usando fetch)
            api.fetchUpload(formData)
            .then(response => {
                if (response.ok) {
                    alert("Arquivo enviado com sucesso! /");
                } else {
                    alert("Erro ao enviar o arquivo.");
                }
            })
            .catch(error => {
                console.error('Erro:', error);
                alert("Erro ao enviar o arquivo.");
            });
        };
    },
    openTerminal(containerId){
        const term = new Terminal();
        const elmTerminal = document.getElementById('terminal-container')
        // clear terminal old
        elmTerminal.innerHTML = ''
        term.open(elmTerminal);
        term.clear()

        // Conecta ao WebSocket do backend
        const socket = api.socket(containerId);
        socket.binaryType = "arraybuffer";

        // Variáveis para o histórico de comandos e buffer de entrada
        let inputBuffer = '';
        let commandHistory = [];
        let historyIndex = -1;

        // Quando o WebSocket recebe uma mensagem do backend
        socket.onmessage = function(event) {
            const message = new Uint8Array(event.data); // Converte o array binário para Uint8Array
            term.write(message); // Escreve a saída do contêiner no terminal
        };

        // Função para substituir o input atual
        function replaceInput(newInput) {
            // Apaga a linha atual
            term.write('\x1b[2K\r'); // ESC sequence para limpar a linha e voltar ao início
            term.write('# '+newInput); // Escreve o novo input
            inputBuffer = newInput; // Atualiza o buffer de entrada
        }

        // Quando o usuário digita no terminal
        term.onData(function(input) {
            
            if (input === '\r') {
                // Enter: envia o comando ao WebSocket e adiciona ao histórico
                socket.send(inputBuffer + '\n');
                commandHistory.push(inputBuffer); // Salva o comando no histórico
                historyIndex = -1; // Reseta o índice do histórico
                inputBuffer = '';
                term.write('\r\n');
            } else if (input === '\u001b[A') { // Seta para cima
                // Recupera o comando anterior do histórico
                if (commandHistory.length > 0 && historyIndex < commandHistory.length - 1) {
                    historyIndex++;
                    replaceInput(commandHistory[commandHistory.length - 1 - historyIndex]);
                }
            } else if (input === '\u001b[B') { // Seta para baixo
                // Move para frente no histórico
                if (historyIndex > 0) {
                    historyIndex--;
                    replaceInput(commandHistory[commandHistory.length - 1 - historyIndex]);
                } else if (historyIndex === 0) {
                    historyIndex = -1;
                    replaceInput(''); // Limpa a linha
                }
            } else if (input === '\x7f') { // Backspace
                // Remove o último caractere do inputBuffer
                if (inputBuffer.length > 0) {
                    inputBuffer = inputBuffer.slice(0, -1);
                    // Move o cursor para trás, apaga o caractere e volta o cursor
                    term.write('\b \b');
                }
            } else if (input === '\x03') { // Ctrl+C
                // Envia o sinal de interrupção (Ctrl+C) para o WebSocket
                socket.send('\x03');
                term.write('^C\r\n'); // Exibe a interrupção no terminal
                inputBuffer = ''; // Limpa o buffer de input
            }  else if (input === '\x18') { // Ctrl+X
                // Envia o sinal de fim de arquivo ou saída (Ctrl+X) para o WebSocket
                socket.send('\x18');
                term.write('^X\r\n'); // Exibe a interrupção no terminal
                inputBuffer = ''; // Limpa o buffer de input
            } else {
                // Captura todos os outros inputs
                inputBuffer += input;
                term.write(input);
            }
        });

        // Detecta quando o WebSocket é desconectado
        socket.onclose = function() {
            term.write('\r\n*** Disconnected from server ***\r\n');
        };

        socket.onerror = function(error) {
            term.write(`\r\n*** Error: ${error.message} ***\r\n`);
            socket.close()
        };
    },
    util:{
        formatBytes(bytes) {
            const sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB'];
            if (bytes === 0) return '0 Byte';
            const i = parseInt(Math.floor(Math.log(bytes) / Math.log(1024)));
            return Math.round(bytes / Math.pow(1024, i), 2) + ' ' + sizes[i];
        }
    }   
}
