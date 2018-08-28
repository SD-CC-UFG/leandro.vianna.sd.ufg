package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

const PORT = 7777

func readMessages() {
	conn, err := net.Dial("tcp", fmt.Sprintf("127.0.0.1:%d", PORT))

	if err != nil {
		log.Fatal("ERROR: %v\n", err)
		return
	}

	defer conn.Close()

	buffer := make([]byte, 4096)

	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	writer.Write([]byte("VIEW\n"))
	writer.Flush()

	for {
		n, err := reader.Read(buffer)

		if err != nil {
			log.Printf("ERROR: %v\n", err)
			break
		}

		response := string(buffer[:n])

		if response == "ERR" {
			log.Printf("ERROR: Server response with error.\n")
			return
		}

		fmt.Println(string(buffer[:n]))
	}
}

func writeMessages(name string, exitChan chan bool) {
	defer func() {
		exitChan <- true
	}()

	conn, err := net.Dial("tcp", fmt.Sprintf("127.0.0.1:%d", PORT))

	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
		return
	}

	defer conn.Close()

	buffer := make([]byte, 4096)

	writer := bufio.NewWriter(conn)
	reader := bufio.NewReader(conn)

	writer.Write([]byte(fmt.Sprintf("JOIN %s\n", name)))
	writer.Flush()

	var msg string

	scanner := bufio.NewScanner(os.Stdin)

	for {
		n, err := reader.Read(buffer)

		if err != nil {
			log.Fatal("Erro fatal: %v\n", err)
		}

		response := string(buffer[:n])

		if response == "ERR" {
			log.Fatal("Erro fatal: Servidor respondeu com erro.\n")
		}

		scanner.Scan()
		msg = scanner.Text()

		if err := scanner.Err(); err != nil {
			log.Fatal("Erro fatal: %v\n", err)
		}

		if msg == "quit" {
			fmt.Printf("Tchau!\n")
			break
		}

		writer.Write([]byte(fmt.Sprintf("MSG %s <end>\n", msg)))
		writer.Flush()
	}
}

func main() {
	var name string
	for len(name) == 0 {
		fmt.Printf("Bem vindo. Digite o seu nome: ")
		fmt.Scan(&name)
	}

	exitChan := make(chan bool)

	go writeMessages(name, exitChan)
	go readMessages()

	<-exitChan
}
