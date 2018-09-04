package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
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

func listenClients(port int) error {
	log.Printf("Esperando conexões na porta %d\n", port)

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

	// Comandos enviados pelo Cliente
	// JOIN username\n
	// MSG message...\n
	// VIEW\n

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
		log.Printf("Nova conexão do tipo JOIN requisitada.\n")

		scanner.Scan()
		name := scanner.Text()

		err = scanner.Err()

		if err != nil {
			log.Printf("ERRO em conexão com cliente: %v", err)
			writer.Write([]byte("ERR\n"))
			writer.Flush()
			break
		}

		if len(name) == 0 {
			log.Printf("Usuário enviou o nome vazio. Encerrando conexão.\n")
			writer.Write([]byte("ERR\n"))
			writer.Flush()
			break
		}

		writer.Write([]byte("OK\n"))
		writer.Flush()
		log.Printf("Nova conexão do tipo JOIN feita por %s\n", name)

		msg := fmt.Sprintf("%s entrou no chat.", name)

		AddMessage("ADM", msg)

		// loop para esperar por mensagens desse usuario
		for {
			scanner.Scan()
			command = scanner.Text()

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

			// adiciona mensagem na fila e volta a esperar novas mensagens
			AddMessage(name, msg)

			writer.Write([]byte("OK\n"))
			writer.Flush()
		}

		msg = fmt.Sprintf("%s saiu do chat.", name)
		AddMessage("ADM", msg)

		log.Printf("Fim de conexao com %s\n", name)

	case "VIEW":
		log.Printf("Nova conexao do tipo VIEW\n")
		var i = 0

		for {
			if !connectAlive(conn) {
				break
			}

			// se houver novas mensagens, envia para o cliente
			for i < len(messageQueue) {
				writer.Write([]byte(
					fmt.Sprintf("%v\n", messageQueue[i])))
				i++
			}
			writer.Flush()
		}
	}

}

func connectAlive(conn net.Conn) bool {
	conn.SetReadDeadline(time.Now())
	var one []byte

	if _, err := conn.Read(one); err == io.EOF {
		return false
	} else {
		var zero time.Time
		conn.SetReadDeadline(zero)
		return true
	}
}

func main() {
	const PORT = 7777
	err := listenClients(PORT)
	if err != nil {
		log.Fatal("Erro fatal: %v\n", err)
	}
}
