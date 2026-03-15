package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/ghibranalj/janus-ssh/tui"
)

// ServerRepository handles server persistence
type ServerRepository struct {
	path string
}

func NewServerRepository(path string) *ServerRepository {
	return &ServerRepository{path: path}
}

func (r *ServerRepository) List() ([]tui.Server, error) {
	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(r.path), 0755); err != nil {
		return nil, fmt.Errorf("failed to create config directory: %w", err)
	}

	data, err := os.ReadFile(r.path)
	if err != nil {
		if os.IsNotExist(err) {
			return []tui.Server{}, nil
		}
		return nil, fmt.Errorf("failed to read servers file: %w", err)
	}

	var servers []tui.Server
	if err := json.Unmarshal(data, &servers); err != nil {
		return nil, fmt.Errorf("failed to parse servers file: %w", err)
	}

	return servers, nil
}

func (r *ServerRepository) Add(server tui.Server) error {

	servers, err := r.List()
	if err != nil {
		return err
	}

	// Check for duplicates
	for _, s := range servers {
		if s.User == server.User && s.Host == server.Host {
			return nil // Already exists
		}
	}

	servers = append(servers, server)

	return r.save(servers)
}

func (r *ServerRepository) Update(oldUser, oldHost string, newServer tui.Server) error {
	servers, err := r.List()
	if err != nil {
		return err
	}

	for i, s := range servers {
		if s.User == oldUser && s.Host == oldHost {
			servers[i] = newServer
			return r.save(servers)
		}
	}

	return fmt.Errorf("server not found")
}

func (r *ServerRepository) Delete(user, host string) error {
	servers, err := r.List()
	if err != nil {
		return err
	}

	var newServers []tui.Server
	for _, s := range servers {
		if s.User != user || s.Host != host {
			newServers = append(newServers, s)
		}
	}

	return r.save(newServers)
}

func (r *ServerRepository) save(servers []tui.Server) error {
	if err := os.MkdirAll(filepath.Dir(r.path), 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	data, err := json.MarshalIndent(servers, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal servers: %w", err)
	}

	if err := os.WriteFile(r.path, data, 0644); err != nil {
		return fmt.Errorf("failed to write servers file: %w", err)
	}

	return nil
}
