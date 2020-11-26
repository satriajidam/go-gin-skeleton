package server

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kelseyhightower/envconfig"
	"github.com/satriajidam/go-gin-skeleton/pkg/log"
)

type config struct {
	GinMode                  string `envconfig:"GIN_MODE" default:"release"`
	GinDisallowUnknownFields bool   `envconfig:"GIN_DISALLOW_UNKNOWN_FIELDS" default:"false"`
}

const (
	ginDebugMode   = "debug"
	ginReleaseMode = "release"
	ginTestMode    = "test"
)

// Server is an interface for all type of servers.
type Server interface {
	Start() error
	Stop(ctx context.Context) error
}

func init() {
	serverConfig := &config{}
	envconfig.MustProcess("", serverConfig)

	switch strings.ToLower(serverConfig.GinMode) {
	case ginDebugMode:
		gin.SetMode(gin.DebugMode)
	case ginReleaseMode:
		gin.SetMode(gin.ReleaseMode)
	case ginTestMode:
		gin.SetMode(gin.TestMode)
	default:
		panic(fmt.Errorf("unsupported gin mode: %s", serverConfig.GinMode))
	}

	if serverConfig.GinDisallowUnknownFields {
		gin.EnableJsonDecoderDisallowUnknownFields()
	}
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

// RunServersGracefully runs all given servers in a graceful way.
func RunServersGracefully(timeoutSeconds time.Duration, servers ...Server) {
	if err := <-StartServers(servers...); err != nil {
		panic(err)
	}

	// Graceful shutdown:
	// - https://chenyitian.gitbooks.io/gin-web-framework/docs/38.html
	// - https://medium.com/honestbee-tw-engineer/gracefully-shutdown-in-go-http-server-5f5e6b83da5a
	// Wait for interrupt signal to gracefully shutdown the server.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// Set graceful shutdown timeout to configured seconds.
	ctx, cancel := context.WithTimeout(
		context.Background(),
		timeoutSeconds*time.Second,
	)
	defer cancel()

	StopServers(ctx, servers...)
}
