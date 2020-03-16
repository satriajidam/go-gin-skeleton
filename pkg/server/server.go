package server

// Server is an interface for all type of servers.
type Server interface {
	Start() error
}

// StartServers starts and run all given servers.
func StartServers(servers ...Server) <-chan error {
	serversCount := len(servers)
	ch := make(chan error, serversCount)

	for _, server := range servers {
		go func(server Server) {
			if err := server.Start(); err != nil {
				ch <- err
			}
		}(server)
	}

	return ch
}
