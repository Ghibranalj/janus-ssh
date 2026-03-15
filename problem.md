# Problem: Bubble Tea v2 Shows Blank Output Over SSH Session

## Application Overview

**janus-ssh** is an SSH jump server that presents a TUI menu for selecting SSH servers to connect to. It uses:
- `github.com/gliderlabs/ssh` v0.3.8 - SSH server library
- `charm.land/bubbletea/v2` v2.0.2 - TUI framework

## The Problem

When a client connects via SSH (`ssh -p 2222 localhost`), the Bubble Tea v2 menu shows **completely blank** - no characters, no cursor, nothing visible.

However, the TUI displays correctly when output goes to stdout (tested by running the application directly in terminal).

## Current Implementation

### File: `/home/gibi/Workspace/janus-ssh/tui.go`

**Key parts:**

```go
// Model struct
type model struct {
    cursor         int
    s              ssh.Session
    selectedServer string
    menu           menu
    username       string
    host           string
    activeField    int
    cursorPos      int
    errorMessage   string
}

// View method - returns tea.View with AltScreen enabled
func (m *model) View() tea.View {
    var view tea.View
    switch m.menu {
    case selectMenu:
        view = tea.NewView(m.selectMenu())
    case addMenu:
        view = tea.NewView(m.addMenu())
    default:
        view = tea.NewView("")
    }
    view.AltScreen = true
    return view
}

// runMenu function - creates Bubble Tea program
func runMenu(s ssh.Session) (string, error) {
    p := tea.NewProgram(
        &model{
            s:    s,
            menu: selectMenu,
        },
        tea.WithInput(s),
        tea.WithOutput(s),
    )

    finalModel, err := p.Run()
    if err != nil {
        return "", err
    }

    m := finalModel.(*model)
    return m.selectedServer, nil
}
```

### File: `/home/gibi/Workspace/janus-ssh/ssh.go`

```go
func SshHandler(s ssh.Session) {
    ptyReq, winCh, ok := s.Pty()

    // restoreTerminal(s) - currently commented out for testing

    for {
        selectedServer, err := runMenu(s)
        if err != nil {
            io.WriteString(s, fmt.Sprintf("Error: %v\n", err))
            return
        }
        restoreTerminal(s)
        if selectedServer == "" {
            return
        }
        // ... rest of SSH handling
    }
}
```

## What Has Been Tried

1. ✅ Set `view.AltScreen = true` in the View() method
2. ✅ Commented out `restoreTerminal()` calls (was suspected to interfere)
3. ✅ Verified the Add Menu feature logic is correctly implemented

## Key Observations

1. **Works on stdout** - Running the code with stdout as output shows the menu correctly
2. **Blank on SSH session** - Using `tea.WithOutput(s)` where `s` is `ssh.Session` results in blank output
3. **Completely blank** - No characters, no cursor, nothing appears in the SSH session
4. **No errors** - `p.Run()` returns without error, but nothing is displayed

## Technical Context

### Bubble Tea v2 API

Bubble Tea v2 uses a `tea.View` struct instead of string:
```go
type View struct {
    Content string
    AltScreen bool
    Cursor *Cursor
    // ... other fields
}
```

The `AltScreen` field enables alternate screen buffer (full window mode).

### gliderlabs/ssh Session

The `ssh.Session` from gliderlabs/ssh implements `io.Reader` and `io.Writer`, which should be compatible with `tea.WithInput()` and `tea.WithOutput()`.

### PTY Support

The code does capture PTY information:
```go
ptyReq, winCh, ok := s.Pty()
```

But this information is not currently passed to `runMenu()` or used in Bubble Tea initialization.

## Expected Behavior

When connecting via `ssh -p 2222 localhost`, the user should see:
```
Select server:

> gibi@reverse-proxy.infra
  gibi@dns.infra

q to quit | a: Add new server
```

## Actual Behavior

The SSH session connects but shows a completely blank screen with no visible content.

## Possible Root Causes

1. **Terminal capability detection** - Bubble Tea v2 may detect the SSH session as a non-terminal
2. **Missing terminal initialization** - SSH sessions may need explicit terminal mode setup (raw mode, echo off)
3. **Interface compatibility** - The SSH session may not implement an interface Bubble Tea expects
4. **Buffering issue** - Output may be buffered but not flushed to the SSH session
5. **PTY not used** - Bubble Tea may need the actual PTY file descriptor, not the SSH session wrapper

## Relevant File Paths

- `/home/gibi/Workspace/janus-ssh/tui.go` - TUI implementation
- `/home/gibi/Workspace/janus-ssh/ssh.go` - SSH handler
- `/home/gibi/Workspace/janus-ssh/main.go` - Entry point
- `/home/gibi/Workspace/janus-ssh/go.mod` - Dependencies

## Dependencies

```
require (
    charm.land/bubbletea/v2 v2.0.2
    github.com/creack/pty v1.1.24
    github.com/gliderlabs/ssh v0.3.8
)
```

## Goal

Fix the blank output issue so that the Bubble Tea v2 TUI displays correctly when rendered over an SSH session using gliderlabs/ssh.
