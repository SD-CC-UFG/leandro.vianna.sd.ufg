package main

import (
	"github.com/sd-cc-ufg/leandro.vianna.sd.ufg/chat_server/server"
	"log"
	"net"
	"strconv"
	"sync"
)

var available = 0
var mutex *sync.Mutex

func Dispatcher(port int, numberThreads int, minAvailable int) error {
	connChannel := make(chan net.Conn)
	mutex = &sync.Mutex{}

	for i := 0; i < numberThreads; i++ {
		// iniciando goroutine e passando o channel para
		// comunicao com dispatcher
		go handler(connChannel, i)
	}

	listen, err := net.Listen("tcp", ":"+strconv.Itoa(port))

	if err != nil {
		return err
	}

	defer listen.Close()

	mutex.Lock()
	available = numberThreads
	mutex.Unlock()

	for {
		conn, err := listen.Accept()
		if err != nil {
			return err
		}

		// enviando conexao para o channel
		// assim a primeira goroutine disponivel vai assumir
		// a conexao
		connChannel <- conn

		mutex.Lock()
		if available < minAvailable {
			log.Printf("Dispatcher criando mais goroutines (subindo para %d)\n", 2*numberThreads)
			for i := numberThreads; i < 2*numberThreads; i++ {
				go handler(connChannel, i)
			}
			available += numberThreads
			numberThreads = 2 * numberThreads
		}
		mutex.Unlock()
	}
}

func handler(connChannel chan net.Conn, mynumber int) {
	for {
		log.Printf("Goroutine %d esperando por uma conexao\n", mynumber)
		conn := <-connChannel

		log.Printf("Goroutine %d recebeu uma conexao\n", mynumber)
		mutex.Lock()
		available--
		mutex.Unlock()

		// passando para servidor tratar conexao
		// ele deve fecha-la
		server.HandleConnection(conn)

		mutex.Lock()
		available++
		mutex.Unlock()
	}
}
