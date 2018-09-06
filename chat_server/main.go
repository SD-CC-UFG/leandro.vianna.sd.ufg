package main

import (
	"log"
)

func main() {
	const PORT = 7777
	const MAX_THREADS = 7
	const MIN_AVAILABLE = 5

	err := Dispatcher(PORT, MAX_THREADS, MIN_AVAILABLE)

	if err != nil {
		log.Fatal(err)
	}
}
