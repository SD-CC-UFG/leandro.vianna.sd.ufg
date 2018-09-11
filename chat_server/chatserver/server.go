package chatserver

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
)

type ChatServer struct {
	messageChannels      []chan Message
	messageChannelsMutex *sync.Mutex
}

func NewChatServer() *ChatServer {
	return &ChatServer{messageChannels: make([]chan Message, 0), messageChannelsMutex: &sync.Mutex{}}
}

func (cs *ChatServer) addChatMessage(name, msg string) {
	log.Printf("Nova mensagem adicionada na fila por %s\n", name)
	message := NewMessage(name, msg)

	cs.messageChannelsMutex.Lock()

	// percorrendo slice com channels para envio de mensagens
	for _, ch := range cs.messageChannels {
		// envio nao bloqueante para o channel ch
		// caso o channel ch nao tenha ninguem recebendo default é executado
		select {
		case ch <- message:
		default:
		}
	}

	cs.messageChannelsMutex.Unlock()
}

func (cs *ChatServer) addMessageChannel(channel chan Message) {
	cs.messageChannelsMutex.Lock()

	cs.messageChannels = append(cs.messageChannels, channel)

	cs.messageChannelsMutex.Unlock()
}

func (cs *ChatServer) Handle(conn net.Conn) {
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
		cs.talkHandler(scanner, writer)

	case "VIEW":
		cs.viewHandler(writer)
	}

	log.Printf("Fim de conexao com %s.\n", conn.RemoteAddr())
}

func (cs *ChatServer) talkHandler(scanner *bufio.Scanner, writer *bufio.Writer) {
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

	cs.addChatMessage("ADM", msg)

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
		cs.addChatMessage(name, msg)

		writer.Write([]byte("OK\n"))
		writer.Flush()
	}

	msg = fmt.Sprintf("%s saiu do chat.", name)
	cs.addChatMessage("ADM", msg)

	log.Printf("Fim de conexao com %s\n", name)
}

func (cs *ChatServer) viewHandler(writer *bufio.Writer) {
	// envia mensagens do chat para cliente conectado

	log.Printf("Nova conexao do tipo VIEW\n")

	mychannel := make(chan Message)
	// limitação: sempre adiciona novos channels ao slice, mantendo channels possivelmente "mortos"
	cs.addMessageChannel(mychannel)

	for {
		// se houver novas mensagens, envia para o cliente
		msg := <-mychannel

		_, err := writer.Write([]byte(fmt.Sprintf("%v\n", msg)))
		if err != nil {
			log.Printf("ERRO em conexao VIEW: %v.", err)
			return
		}
		writer.Flush()
	}
}
