package main

import (
	"bufio"
	"log"
	"net"
	"strconv"
	"strings"
	"time"
)

func dispatcher(port int, numberThreads int) error {
	connChanQueue := make([]chan net.Conn, numberThreads)

	for i := 0; i < numberThreads; i++ {
		connChanQueue[i] = make(chan net.Conn)
		go handle(connChanQueue[i], i)
	}

	listen, err := net.Listen("tcp", ":"+strconv.Itoa(port))

	if err != nil {
		return err
	}

	defer listen.Close()

	var nextGoroutine = 0

	for {
		conn, err := listen.Accept()
		if err != nil {
			return err
		}

		log.Printf("Dispatcher scheduled goroutine %d.\n", nextGoroutine)
		connChanQueue[nextGoroutine] <- conn
		nextGoroutine++
		nextGoroutine %= numberThreads
	}
}

func handle(connChan chan net.Conn, mynumber int) {
	for {
		log.Printf("Goroutine %d expecting connection\n", mynumber)
		conn := <-connChan
		log.Printf("Goroutine %d accepted a connection\n", mynumber)

		reader := bufio.NewReader(conn)
		writer := bufio.NewWriter(conn)

		name, err := reader.ReadString('\n')

		if err != nil {
			log.Printf("Error in connection with %s: %v\n", conn.RemoteAddr(), err)
			break
		}

		duration, err := time.ParseDuration("10s")
		time.Sleep(duration)

		name = strings.ToUpper(name)
		writer.Write([]byte(name))
		writer.Flush()
		conn.Close()
	}
}

func main() {
	const PORT = 7777

	err := dispatcher(PORT, 10)

	if err != nil {
		log.Fatal(err)
	}
}
