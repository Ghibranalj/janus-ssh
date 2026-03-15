package tui

type Server struct {
	User string
	Host string
}

// String returns the server in "user@host" format
func (s Server) String() string {
	return s.User + "@" + s.Host
}

// ServerRepository interface allows tui to persist servers without importing main
type ServerRepository interface {
	List() ([]Server, error)
	Add(server Server) error
}
