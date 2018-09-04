package main

import (
	"errors"
	"github.com/sd-cc-ufg/leandro.vianna.sd.ufg/chat_server/queue"
	"github.com/sd-cc-ufg/leandro.vianna.sd.ufg/chat_server/server"
	"log"
	"net"
	"strconv"
)

var connChanQueue *queue.Queue

func Dispatcher(port int, numberThreads int) error {
	// criando fila com channels de conexoes
	connChanQueue = queue.New()

	for i := 0; i < numberThreads; i++ {
		// iniciando goroutine e armazendo na fila channel para
		// comunicao com dispatcher
		channel := make(chan net.Conn) // canal de comunicao dispatcher <-> goroutine
		connChanQueue.Push(channel)
		go handler(channel, i)
	}

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

		if !connChanQueue.Empty() {
			channel, err := connChanQueue.Pop()
			if err != nil {
				conn.Close()
				return err
			}

			switch channel := channel.(type) {
			case chan net.Conn:
				channel <- conn
			default:
				conn.Close()
				return errors.New("Item of channels queue is not net.Conn channel\n")
			}
		}
	}
}

func handler(connChan chan net.Conn, mynumber int) {
	for {
		log.Printf("Goroutine %d expecting connection\n", mynumber)
		conn := <-connChan
		log.Printf("Goroutine %d accepted a connection\n", mynumber)

		// passando para servidor tratar conexao
		// ele deve fecha-la
		server.HandleConnection(conn)

		// ao final do tratamento
		// goroutine coloca seu channel na fila para ficar disponivel
		connChanQueue.Push(connChan)
	}
}
