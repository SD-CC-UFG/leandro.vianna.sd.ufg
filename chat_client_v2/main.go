package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

const PORT = 7777

func writeMessages(name string) {
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

	// verifica status da conexao retornado pelo server
	n, err := reader.Read(buffer)

	if err != nil {
		log.Fatalf("Erro fatal: %v\n", err)
	}

	response := string(buffer[:n])

	if response == "ERR" {
		log.Fatal("Erro fatal: Servidor respondeu com erro.\n")
		return
	}

	go func() {
		for {
			n, err := reader.Read(buffer)

			if err != nil {
				log.Fatalf("ERROR: %v\n", err)
			}

			response := string(buffer[:n])

			if response == "ERR" {
				log.Fatal("ERROR: Server response with error.\n")
			}

			fmt.Print(string(buffer[:n]))
		}
	}()

	stdin := bufio.NewScanner(os.Stdin)

	for {
		// recebe texto do usuario
		stdin.Scan()
		msg = stdin.Text()

		if err := stdin.Err(); err != nil {
			log.Fatalf("Erro fatal: %v\n", err)
		}

		if msg == "quit" {
			writer.Write([]byte("quit"))
			writer.Flush()
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

	writeMessages(name)
}
