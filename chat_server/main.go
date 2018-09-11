package main

import (
	"github.com/sd-cc-ufg/leandro.vianna.sd.ufg/chat_server/chatserver"
	"github.com/sd-cc-ufg/leandro.vianna.sd.ufg/chat_server/dispatcher"
	"log"
)

func main() {
	const PORT = 7777
	const MAX_THREADS = 10
	const MIN_AVAILABLE = 4

	chatServer := chatserver.NewChatServer()
	dispatcher := dispatcher.NewDispatcher(PORT, MAX_THREADS, MIN_AVAILABLE, chatServer)

	err := dispatcher.Start()

	if err != nil {
		log.Fatal(err)
	}
}
