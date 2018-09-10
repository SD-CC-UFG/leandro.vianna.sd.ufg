package main

import (
	"github.com/sd-cc-ufg/leandro.vianna.sd.ufg/chat_server/dispatcher"
	"github.com/sd-cc-ufg/leandro.vianna.sd.ufg/chat_server/server"
	"log"
)

func main() {
	const PORT = 7777
	const MAX_THREADS = 10000
	const MIN_AVAILABLE = 10

	dispatcher := dispatcher.NewDispatcher(PORT, MAX_THREADS, MIN_AVAILABLE,
		server.HandleConnection)

	err := dispatcher.Start()

	if err != nil {
		log.Fatal(err)
	}
}
