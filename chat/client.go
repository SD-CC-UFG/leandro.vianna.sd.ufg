package main

import (
	"fmt"
	"net"
)

const PORT = 7777

func readMessages() {
	conn, err := net.Dial("tcp", fmt.Sprintf("localhost:%d", PORT))

	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
		return
	}

	defer conn.Close()

	buffer := make([]byte, 4096)

	conn.Write([]byte("VIEW"))

	for {
		n, err := conn.Read(buffer)

		if err != nil {
			fmt.Printf("ERROR: %v\n", err)
			break
		}

		fmt.Println(string(buffer[:n]))
	}
}

func writeMessages(name string) {
	var msg string

	conn, err := net.Dial("tcp", fmt.Sprintf("127.0.0.1:%d", PORT))

	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
		return
	}

	defer conn.Close()

	conn.Write([]byte("JOIN"))

	conn.Write([]byte(name))
	fmt.Printf("Send name and JOIN\n")

	for {
		fmt.Scanln(&msg)
		conn.Write([]byte(msg))

		if msg == "quit" {
			break
		}
	}
}

func main() {
	var name string
	for len(name) == 0 {
		fmt.Printf("Bem vindo. Digite o seu nome: ")
		fmt.Scan(&name)
	}

	go writeMessages(name)
	readMessages()
}
