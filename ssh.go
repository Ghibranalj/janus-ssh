package main

import (
	"context"
	"fmt"
	"os"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/wish/v2"
	"charm.land/wish/v2/activeterm"
	"charm.land/wish/v2/bubbletea"
	"charm.land/wish/v2/logging"
	"github.com/charmbracelet/ssh"

	"github.com/ghibranalj/janus-ssh/tui"
)

type SSHServer struct {
	address     string
	hostKeyPath string
	repo        tui.ServerRepository
	server      *ssh.Server
}

func NewSSHServer(address, hostKeyPath string, repo tui.ServerRepository) *SSHServer {
	return &SSHServer{
		address:     address,
		hostKeyPath: hostKeyPath,
		repo:        repo,
	}
}

func (s *SSHServer) Start() error {
	var err error
	s.server, err = wish.NewServer(
		wish.WithAddress(s.address),
		wish.WithMiddleware(
			bubbletea.Middleware(s.teaHandler),
			activeterm.Middleware(), // Bubble Tea apps usually require a PTY.
			logging.Middleware(),
		),
	)
	if err != nil {
		return fmt.Errorf("could not create server: %w", err)
	}

	done := make(chan os.Signal, 1)

	fmt.Printf("Starting SSH server on %s\n", s.address)
	go func() {
		if err := s.server.ListenAndServe(); err != nil {
			fmt.Fprintf(os.Stderr, "Server error: %v\n", err)
		}
	}()

	<-done
	fmt.Println("Shutting down server...")
	return s.Shutdown()
}

func (s *SSHServer) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := s.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("could not shutdown server: %w", err)
	}
	return nil
}

func (s *SSHServer) teaHandler(sess ssh.Session) (tea.Model, []tea.ProgramOption) {
	return tui.NewApp(s.repo, sess), []tea.ProgramOption{}
}
