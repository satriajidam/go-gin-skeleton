package main

import (
	"github.com/satriajidam/go-gin-skeleton/pkg/config"
	"github.com/satriajidam/go-gin-skeleton/pkg/database/sql"
	"github.com/satriajidam/go-gin-skeleton/pkg/database/sql/sqlite"
	"github.com/satriajidam/go-gin-skeleton/pkg/server"
	"github.com/satriajidam/go-gin-skeleton/pkg/server/http"
	"github.com/satriajidam/go-gin-skeleton/pkg/server/prometheus"
	"github.com/satriajidam/go-gin-skeleton/pkg/service/api"
	"github.com/satriajidam/go-gin-skeleton/pkg/service/client/http/pokeapi"
	"github.com/satriajidam/go-gin-skeleton/pkg/service/pokemon"
	"github.com/satriajidam/go-gin-skeleton/pkg/service/provider"
)

func main() {
	cfg := config.Get()

	dbconn, err := sqlite.NewConnection(sql.DBConfig{
		Database:      cfg.SQLiteDatabase,
		MaxIdleConns:  cfg.SQLiteMaxIdleConns,
		MaxOpenConns:  cfg.SQLiteMaxOpenConns,
		SingularTable: cfg.SQLiteSingularTable,
		DebugMode:     cfg.SQLiteDebugMode,
	})
	defer func() {
		err := dbconn.Close()
		if err != nil {
			panic(err)
		}
	}()
	if err != nil {
		panic(err)
	}

	httpServer := http.NewServer(cfg.HTTPServerPort, true)

	providerRepository := provider.NewRepository(dbconn, true)
	providerService := provider.NewService(providerRepository)
	providerHTTPHandler := api.NewProviderHTTPHandler(providerService)

	pokeapiClient := pokeapi.NewClient(cfg.PokeAPIAddressV2, 15)
	pokemonService := pokemon.NewService(pokeapiClient)
	pokemonHTTPHandler := api.NewPokemonHTTPHandler(pokemonService)

	v1 := httpServer.Group("/v1")
	v1.POST("/provider", providerHTTPHandler.CreateProvider)
	v1.PUT("/provider/:uuid", providerHTTPHandler.UpdateProvider)
	v1.DELETE("/provider/:uuid", providerHTTPHandler.DeleteProviderByUUID)
	v1.GET("/provider/:uuid", providerHTTPHandler.GetProviderByUUID)
	v1.GET("/providers", providerHTTPHandler.ListProviders)
	v1.GET("/pokemon/:name", pokemonHTTPHandler.GetPokemonByName)

	skipMonitoringPaths := []string{"/_/health"}

	httpServer.LoggerSkipPaths(skipMonitoringPaths...)

	promServer := prometheus.NewServer(
		cfg.PrometheusServerPort,
		cfg.PrometheusServerMetricsPath,
	)

	promServer.Monitor(
		&prometheus.Target{
			HTTPServer:   httpServer,
			ExcludePaths: skipMonitoringPaths,
		},
	)

	server.RunServersGracefully(cfg.GracefulTimeout, promServer, httpServer)
}
