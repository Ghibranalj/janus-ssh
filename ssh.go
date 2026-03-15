package main

import (
	"fmt"
	"io"
	"os/exec"

	"github.com/creack/pty"
	"github.com/gliderlabs/ssh"

	"github.com/ghibranalj/janus-ssh/tui"
)

func SshHandler(s ssh.Session, repo tui.ServerRepository) {
	ptyReq, winCh, ok := s.Pty()

	tui.RestoreTerminal(s)

	for {
		selectedServer, err := tui.RunMenu(s, repo)
		if err != nil {
			io.WriteString(s, fmt.Sprintf("Error: %v\n", err))
			return
		}
		tui.RestoreTerminal(s)
		if selectedServer == (tui.Server{}) {
			return
		}

		if ok {
			err = runSSHSession(s, selectedServer.String(), ptyReq, winCh)
			if err != nil {
				io.WriteString(s, fmt.Sprintf("SSH error: %v\n", err))
			}
			tui.RestoreTerminal(s)
		} else {
			cmd := exec.Command("ssh", selectedServer.String())
			cmd.Stdin = s
			cmd.Stdout = s
			cmd.Stderr = s
			cmd.Run()
		}
	}
}

// runSSHSession starts the SSH command with PTY support
func runSSHSession(s ssh.Session, server string, ptyReq ssh.Pty, winCh <-chan ssh.Window) error {
	cmd := exec.Command("ssh", server)
	cmd.Env = append(cmd.Env, fmt.Sprintf("TERM=%s", ptyReq.Term))

	ptmx, err := pty.Start(cmd)
	if err != nil {
		return err
	}
	defer ptmx.Close()

	// Window resize handler
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

	// Copy stdin to PTY and stdout from PTY (bidirectional)
	go func() {
		io.Copy(ptmx, s)
	}()
	go func() {
		io.Copy(s, ptmx)
	}()

	return cmd.Wait()
}
