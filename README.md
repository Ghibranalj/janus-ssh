# Janus SSH

> **Janus** is the Roman god of gateways, portals, passages, endings, and time. Often depicted with two faces looking in opposite directions, he guards thresholds and transitions - the perfect metaphor for an SSH bastion that serves as the gateway to your servers.

An interactive SSH bastion server built with Go and [Bubble Tea](https://github.com/charmbracelet/bubbletea). Janus SSH provides a TUI-based interface to manage and quickly connect to your SSH servers without remembering complex addresses.

## Features

- **Interactive TUI** - Beautiful terminal interface for browsing and selecting servers
- **Persistent Server List** - Store and manage your frequently used SSH connections
- **Add/Edit/Delete Servers** - Full CRUD operations directly in the SSH session
- **Bastion Host Pattern** - Single entry point to multiple servers
- **Flexible Auth** - Support for both password and public key authentication
- **Docker Ready** - Deploy anywhere with Docker or Docker Compose

## How It Works

1. Connect to Janus SSH via SSH client
2. Browse your saved servers in the interactive menu
3. Select a server to connect - Janus proxies the SSH connection
4. Manage your server list without leaving the SSH session

## Installation

### From Source

```bash
git clone https://github.com/ghibranalj/janus-ssh.git
cd janus-ssh
go mod download
go build
```

### With Docker

```bash
docker build -t janus-ssh .
```

## Configuration

Janus SSH uses environment variables for configuration:

| Variable | Description | Default |
|----------|-------------|---------|
| `PORT` | SSH port to listen on | `2222` |
| `HOST_KEY_FILE` | Path to host private key | Required |
| `AUTH_KEY_FILE` | Path to authorized_keys file | Optional* |
| `PASSWORD_HASH` | Bcrypt password hash | Optional* |

*At least one of `AUTH_KEY_FILE` or `PASSWORD_HASH` must be set.

### Generate Password Hash

```bash
# Using htpasswd
htpasswd -bnBC 10 "" "yourpassword" | tr -d ':\n'

# Or using Python
python3 -c 'import bcrypt; print(bcrypt.hashpw(b"yourpassword", bcrypt.gensalt(rounds=10)).decode())'
```

### Generate Host Keys

```bash
ssh-keygen -t ed25519 -f ./host_key -N ""
```

## Usage

### Running Directly

```bash
export PORT=2222
export HOST_KEY_FILE=./host_key
export PASSWORD_HASH="<your-bcrypt-hash>"
./janus-ssh
```

### Running with Docker Compose

```bash
# 1. Create directories
mkdir -p keys data

# 2. Generate password hash and add to docker-compose.yml

# 3. Start the service
docker-compose up -d

# 4. Connect
ssh -p 2222 user@localhost
```

### Using Environment File

Create a `.env` file:

```env
PORT=2222
HOST_KEY_FILE=./keys/host_key
PASSWORD_HASH=$2b$10$...
AUTH_KEY_FILE=./keys/authorized_keys
```

Then run:

```bash
docker-compose up -d
```

## TUI Controls

| Key | Action |
|-----|--------|
| `↑`/`↓` or `j`/`k` | Navigate menu |
| `Enter` | Select server / Confirm |
| `a` | Add new server |
| `e` | Edit selected server |
| `d` | Delete selected server |
| `q` or `Ctrl+C` | Quit |

## Project Structure

```
janus-ssh/
├── main.go           # Entry point and config
├── ssh.go            # SSH server implementation
├── persistence.go    # Server data persistence
├── tui/              # Terminal UI components
│   ├── tui.go        # Main TUI app
│   ├── selectMenu.go # Server selection
│   ├── serverForm.go # Add/Edit forms
│   └── styles.go     # Styling constants
├── Dockerfile
├── docker-compose.yml
└── README.md
```

## License

To be determined
