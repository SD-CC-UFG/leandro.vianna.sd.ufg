package main

import (
	"log"
)

func main() {
	const PORT = 7777
	const MAX_THREADS = 1000

	err := Dispatcher(PORT, MAX_THREADS)

	if err != nil {
		log.Fatal(err)
	}
}
