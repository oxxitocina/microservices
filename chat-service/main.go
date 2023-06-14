package main

import (
	"log"
)

func main() {
	go startServer()
	startClient()
}

func startServer() {
	if err := server.Run(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func startClient() {
	if err := client.Run(); err != nil {
		log.Fatalf("Failed to start client: %v", err)
	}
}
