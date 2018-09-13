package dispatcher

import (
	"log"
	"net"
	"strconv"
	"sync"
)

type ConnectionHandler interface {
	Handle(conn net.Conn)
}

type Dispatcher struct {
	available int
	mutex     *sync.Mutex

	port          int
	numberThreads int
	minAvailable  int

	handler ConnectionHandler
}

func NewDispatcher(port, numberThreads, minAvailable int, handler ConnectionHandler) Dispatcher {
	dispatcher := Dispatcher{available: 0, mutex: &sync.Mutex{}, port: port,
		numberThreads: numberThreads, minAvailable: minAvailable, handler: handler}
	return dispatcher
}

func (d *Dispatcher) Start() error {
	connChannel := make(chan net.Conn)

	for i := 0; i < d.numberThreads; i++ {
		// iniciando goroutine e passando o channel para
		// comunicao com dispatcher
		go d.waitPassConnection(connChannel, i)
	}

	listen, err := net.Listen("tcp", ":"+strconv.Itoa(d.port))

	if err != nil {
		return err
	}

	defer listen.Close()

	d.mutex.Lock()
	d.available = d.numberThreads
	d.mutex.Unlock()

	for {
		conn, err := listen.Accept()
		if err != nil {
			return err
		}

		// enviando conexao para o channel
		// assim a primeira goroutine disponivel vai assumir
		// a conexao
		connChannel <- conn

		d.mutex.Lock()
		if d.available < d.minAvailable {
			log.Printf("Dispatcher criando mais goroutines (subindo para %d)\n", 2*d.numberThreads)
			for i := d.numberThreads; i < 2*d.numberThreads; i++ {
				go d.waitPassConnection(connChannel, i)
			}
			d.available += d.numberThreads
			d.numberThreads = 2 * d.numberThreads
		}
		d.mutex.Unlock()
	}
}

func (d *Dispatcher) waitPassConnection(connChannel chan net.Conn, mynumber int) {
	for {
		log.Printf("Goroutine %d esperando por uma conexao\n", mynumber)
		conn := <-connChannel

		log.Printf("Goroutine %d recebeu uma conexao\n", mynumber)
		d.mutex.Lock()
		d.available--
		d.mutex.Unlock()

		// passando conexao para callback passada na criacao
		// do dispatcher tratar a conexao
		d.handler.Handle(conn)

		d.mutex.Lock()
		d.available++
		d.mutex.Unlock()
	}
}
