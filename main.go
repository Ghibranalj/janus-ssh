package main

import (
	"fmt"
	"log"
)

func main() {
	repo := NewServerRepository("./servers.json")

	server := NewSSHServer(
		"localhost:2222",
		".ssh/id_ed25519",
		repo,
	)

	if err := server.Start(); err != nil {
		log.Fatalf("Server error: %v", err)
	}
	fmt.Println("Server stopped")
}
