package main

import (
	"time"

	"github.com/satriajidam/go-gin-skeleton/pkg/server"
	"github.com/satriajidam/go-gin-skeleton/pkg/server/http"
	"github.com/satriajidam/go-gin-skeleton/pkg/server/prometheus"
)

func main() {
	httpServer := http.NewServer("8080", true, true)
	promServer := prometheus.NewServer("9180", "/metrics")

	promServer.Monitor(
		&prometheus.Target{
			HTTPServer:             httpServer,
			ExcludePaths:           []string{""},
			GroupedStatus:          false,
			DisableMeasureReqSize:  false,
			DisableMeasureRespSize: false,
			DisableMeasureInflight: false,
		},
	)

	server.RunServersGracefully(time.Second*5, promServer, httpServer)
}
