package main

import (
	"github.com/gliderlabs/ssh"
)

func main() {
	// Set up server repository
	repo := NewServerRepository("./servers.json")

	ssh.Handle(func(s ssh.Session) {
		SshHandler(s, repo)
	})

	ssh.ListenAndServe(":2222", nil)
}
