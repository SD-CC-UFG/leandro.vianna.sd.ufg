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
	messageChannels      map[string]chan Message
	messageChannelsMutex *sync.Mutex
}

func NewChatServer() *ChatServer {
	return &ChatServer{messageChannels: map[string]chan Message{},
		messageChannelsMutex: &sync.Mutex{}}
}

func (cs *ChatServer) addChatMessage(name, msg string) {
	log.Printf("Nova mensagem adicionada na fila por %s\n", name)
	message := NewMessage(name, msg)

	cs.messageChannelsMutex.Lock()

	// percorrendo slice com channels para envio de mensagens
	for k, ch := range cs.messageChannels {
		log.Printf("Enviando mensagens para cliente %s.\n", k)
		ch <- message
	}

	cs.messageChannelsMutex.Unlock()
}

func (cs *ChatServer) addMessageChannel(username string, channel chan Message) {
	cs.messageChannelsMutex.Lock()
	cs.messageChannels[username] = channel
	cs.messageChannelsMutex.Unlock()
}

func (cs *ChatServer) delMessageChannel(username string) {
	cs.messageChannelsMutex.Lock()
	delete(cs.messageChannels, username)
	cs.messageChannelsMutex.Unlock()
}

func (cs *ChatServer) Handle(conn net.Conn) {
	// tratador de conexao do servidor deve fechar conexao
	defer conn.Close()

	// Comandos enviados pelo Cliente
	// JOIN username\n
	// MSG message...\n

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

	if command == "JOIN" {
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

		if _, ok := cs.messageChannels[name]; ok {
			log.Printf("Usuário enviou o nome já existente. Encerrando conexão.\n")
			writer.Write([]byte("ERR\n"))
			writer.Flush()
			return
		}

		writer.Write([]byte("OK\n"))
		writer.Flush()
		log.Printf("Conexao JOIN com %s efetivada\n", name)

		// cria channel de mensagens para este cliente
		mychannel := make(chan Message)
		cs.addMessageChannel(name, mychannel)

		// inicia loop com envio de mensagens em outra rotina
		go cs.sendMessagesFor(mychannel, writer)

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
		}

		msg = fmt.Sprintf("%s saiu do chat.", name)
		cs.addChatMessage("ADM", msg)

		// remove cliente do map de channels de mensagens
		cs.delMessageChannel(name)
		// fecha channel para encerrar rotina de envio de mensagens
		close(mychannel)

		log.Printf("Fim de conexao com %s\n", name)
	} else {
		log.Printf("ERRO em conexão com cliente: %v.", err)
	}
}

func (cs *ChatServer) sendMessagesFor(mychannel chan Message, writer *bufio.Writer) {
	// envia mensagens do chat para cliente conectado

	for {
		// se houver novas mensagens, envia para o cliente
		msg, more := <-mychannel

		if !more {
			// cliente desconectado, encerrar rotina
			break
		}

		_, err := writer.Write([]byte(fmt.Sprintf("%v\n", msg)))
		if err != nil {
			log.Printf("ERRO em conexao VIEW: %v.", err)
			break
		}
		writer.Flush()
	}
}
