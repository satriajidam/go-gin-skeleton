package server

// Server ...
type Server interface {
	Start() error
}

// InitServers ...
func InitServers(servers ...Server) error {
	// TODO: Run each server in its own goroutine.
	for _, server := range servers {
		if err := server.Start(); err != nil {
			return err
		}
	}

	return nil
}
