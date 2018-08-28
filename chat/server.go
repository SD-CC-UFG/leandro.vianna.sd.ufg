package main

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	"sync"
)

type Message struct {
	SenderName string
	Text       string
}

func NewMessage(name, msg string) Message {
	message := Message{name, msg}
	return message
}

var messageQueue []Message
var msgQueueMutex sync.Mutex

func AddMessage(name, msg string) {
	msgQueueMutex.Lock()
	log.Printf("Nova mensagem adicionada na fila por %s\n", name)
	messageQueue = append(messageQueue, NewMessage(name, msg))
	msgQueueMutex.Unlock()
}

func listenClients(port int) error {
	log.Printf("Listening for connections in port %d\n", port)

	listen, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		return err
	}
	defer listen.Close()

	for {
		conn, err := listen.Accept()
		if err != nil {
			return err
		}

		go handleConnection(conn)
	}

	return nil
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	buffer := make([]byte, 4096)

	n, err := conn.Read(buffer)

	if err != nil {
		log.Printf("ERRO em conexão com cliente: %v.", err)
		return
	}

	command := string(buffer[:n])
	command = strings.Trim(command, " ")

	switch command {
	case "JOIN":
		log.Printf("Nova conexão do tipo JOIN iniciada.\n")

		n, err = conn.Read(buffer)

		if err != nil {
			log.Printf("ERRO em conexão com cliente: %v", err)
			break
		}

		name := strings.Trim(string(buffer[:n]), " ")

		if len(name) == 0 {
			log.Printf("Usuário enviou o nome vazio. Encerrando conexão.\n")
			break
		}

		log.Printf("Nova conexão do tipo JOIN feita por %s\n", name)

		msg := fmt.Sprintf("%s entrou no chat.", name)

		AddMessage("ADM", msg)

		// loop para esperar por mensagens desse usuario
		for {
			n, err = conn.Read(buffer)
			if err != nil {
				log.Printf("ERRO em conexão com cliente: %v.", err)
				break
			}

			// recebe mensagem e trata espaços em branco
			msg = string(buffer[:n])
			msg = strings.Trim(msg, " ")

			// verifica se é mensagem vazia ou comando quit
			if n == 0 || msg == "quit" {
				// encerra loop para desconectar usuario
				break
			}

			// adiciona mensagem na fila e volta a esperar novas mensagens
			AddMessage(name, msg)
		}

		msg = fmt.Sprintf("%s saiu do chat.", name)
		AddMessage("ADM", msg)

		log.Printf("Fim de conexao com %s\n", name)

	case "VIEW":
		log.Printf("Nova conexao do tipo VIEW\n")
		var i = 0
		var msg string

		for {
			// se houver novas mensagens, envia para o cliente
			for i < len(messageQueue) {
				msgQueueMutex.Lock()
				msg = fmt.Sprintf("%s: %s\n", messageQueue[i].SenderName, messageQueue[i].Text)
				msgQueueMutex.Unlock()
				conn.Write([]byte(msg))
				i++
			}
		}
	}

}

func main() {
	const PORT = 7777
	err := listenClients(PORT)
	if err != nil {
		log.Fatal("ERRO: %v\n", err)
	}
}
