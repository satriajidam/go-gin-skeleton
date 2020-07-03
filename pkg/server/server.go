package server

import (
	"context"

	"github.com/satriajidam/go-gin-skeleton/pkg/log"
)

// Server is an interface for all type of servers.
type Server interface {
	Start() error
	Stop(ctx context.Context) error
}

// StartServers starts all given servers.
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

// StopServers stops all given servers.
func StopServers(ctx context.Context, servers ...Server) {
	log.Info("Shutting down all servers")
	for _, server := range servers {
		if err := server.Stop(ctx); err != nil {
			log.Fatal(err, "Failed shutting down servers")
		}
	}
	log.Info("All servers exited properly")
}
