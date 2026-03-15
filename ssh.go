package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/wish/v2"
	"charm.land/wish/v2/activeterm"
	"charm.land/wish/v2/bubbletea"
	"charm.land/wish/v2/logging"
	"github.com/charmbracelet/ssh"
	"github.com/creack/pty"

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
			activeterm.Middleware(), // Bubble Tea apps usually require a PTY.
			logging.Middleware(),
		),
	)
	if err != nil {
		return fmt.Errorf("could not create server: %w", err)
	}

	s.server.Handle(s.Handle)

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

func (s *SSHServer) Handle(session ssh.Session) {
	opts := bubbletea.MakeOptions(session)
	tui.RestoreTerminal(session)
	for {
		tuiApp := tui.NewApp(s.repo)
		p := tea.NewProgram(tuiApp, opts...)
		m, err := p.Run()
		if err != nil {
			fmt.Fprintf(session, "Error: %s", err.Error())
		}

		model, ok := m.(*tui.AppModel)

		if !ok {
			fmt.Println("WTF")
			return
		}

		if model.Exit {
			tui.RestoreTerminal(session)
			return
		}
		SSH(session, model.Server)
	}
}

func SSH(session ssh.Session, server string) {
	tui.RestoreTerminal(session)

	fmt.Fprintf(session, "ssh to %s....\n", server)

	ptyReq, winCh, _ := session.Pty()

	cmd := exec.Command("ssh", server)
	cmd.Env = append(cmd.Env, fmt.Sprintf("TERM=%s", ptyReq.Term))

	ptmx, err := pty.Start(cmd)
	if err != nil {
		fmt.Fprintf(session, "Error: %s", err.Error())
	}
	defer ptmx.Close()

	// Handle window size changes
	if winCh != nil {
		go func() {
			for win := range winCh {
				pty.Setsize(ptmx, &pty.Winsize{
					Rows: uint16(win.Height),
					Cols: uint16(win.Width),
				})
			}
		}()
	}

	// Copy bidirectional IO
	// Input: user -> SSH (use goroutine)
	go func() {
		io.Copy(ptmx, session)
	}()
	// Output: SSH -> user (blocking)
	io.Copy(session, ptmx)

	cmd.Wait()
	tui.RestoreTerminal(session)
}
