package server

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
	"time"
)

type Message struct {
	SenderName string
	Text       string
	Time       time.Time
}

func (m Message) String() string {
	return fmt.Sprintf("[%s] %s: %s\n", m.Time.Format("15:04:05"), m.SenderName, m.Text)
}

func NewMessage(name, msg string) Message {
	message := Message{name, msg, time.Now()}
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

func HandleConnection(conn net.Conn) {
	// tratador de conexao do servidor deve fechar conexao
	defer conn.Close()

	// Comandos enviados pelo Cliente
	// JOIN username\n
	// MSG message...\n
	// VIEW\n

	// trata mensagens enviadas pelos clientes
	scanner := bufio.NewScanner(conn)
	scanner.Split(bufio.ScanWords)

	writer := bufio.NewWriter(conn)

	scanner.Scan()
	command := scanner.Text()

	err := scanner.Err()

	if err != nil {
		log.Printf("ERRO em conexão com cliente: %v.", err)
		return
	}

	switch command {
	case "JOIN":
		talkHandler(scanner, writer)

	case "VIEW":
		viewHandler(writer)
	}

	log.Printf("Fim de conexao com %s.\n", conn.RemoteAddr())
}

func talkHandler(scanner *bufio.Scanner, writer *bufio.Writer) {
	log.Printf("Nova conexão do tipo JOIN requisitada.\n")

	scanner.Scan()
	name := scanner.Text()

	err := scanner.Err()

	if err != nil {
		log.Printf("ERRO em conexão com cliente: %v", err)
		writer.Write([]byte("ERR\n"))
		writer.Flush()
		return
	}

	if len(name) == 0 {
		log.Printf("Usuário enviou o nome vazio. Encerrando conexão.\n")
		writer.Write([]byte("ERR\n"))
		writer.Flush()
		return
	}

	writer.Write([]byte("OK\n"))
	writer.Flush()
	log.Printf("Conexao JOIN com %s efetivada\n", name)

	msg := fmt.Sprintf("%s entrou no chat.", name)

	AddMessage("ADM", msg)

	// loop para esperar por mensagens desse usuario
	for {
		scanner.Scan()
		command := scanner.Text()

		err = scanner.Err()
		if err != nil {
			log.Printf("ERRO em conexão com cliente: %v.", err)
			writer.Write([]byte("ERR\n"))
			writer.Flush()
			break
		}

		if command != "MSG" {
			log.Printf("ERRO em conexão com cliente: era esperado comando MSG.")
			writer.Write([]byte("ERR\n"))
			writer.Flush()
			break
		}

		// recebe mensagem e trata espaços em branco
		tokens := make([]string, 0)

		for scanner.Scan() {
			if scanner.Text() == "<end>" {
				break
			}

			tokens = append(tokens, scanner.Text())
		}

		msg = strings.Join(tokens, " ")

		if msg == "quit" {
			break
		}

		// adiciona mensagem e volta a esperar novas mensagens
		AddMessage(name, msg)

		writer.Write([]byte("OK\n"))
		writer.Flush()
	}

	msg = fmt.Sprintf("%s saiu do chat.", name)
	AddMessage("ADM", msg)

	log.Printf("Fim de conexao com %s\n", name)
}

func viewHandler(writer *bufio.Writer) {
	log.Printf("Nova conexao do tipo VIEW\n")
	var i = 0

	// envia mensagens do chat para cliente conectado
	for {

		// se houver novas mensagens, envia para o cliente
		for i < len(messageQueue) {
			_, err := writer.Write([]byte(fmt.Sprintf("%v\n", messageQueue[i])))
			if err != nil {
				log.Printf("ERRO em conexao VIEW: %v.", err)
				return
			}
			i++
		}
		writer.Flush()
	}
}
