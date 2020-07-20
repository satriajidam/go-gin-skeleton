package main

import (
	"fmt"

	"github.com/satriajidam/go-gin-skeleton/pkg/cache/redis"
	"github.com/satriajidam/go-gin-skeleton/pkg/config"
	"github.com/satriajidam/go-gin-skeleton/pkg/database/sql"
	"github.com/satriajidam/go-gin-skeleton/pkg/database/sql/mysql"
	"github.com/satriajidam/go-gin-skeleton/pkg/server"
	"github.com/satriajidam/go-gin-skeleton/pkg/server/http"
	"github.com/satriajidam/go-gin-skeleton/pkg/server/prometheus"
	"github.com/satriajidam/go-gin-skeleton/pkg/service/api"
	"github.com/satriajidam/go-gin-skeleton/pkg/service/client/pokeapi"
	"github.com/satriajidam/go-gin-skeleton/pkg/service/pokemon"
	"github.com/satriajidam/go-gin-skeleton/pkg/service/provider"
)

func main() {
	cfg := config.Get()

	dbconn, err := mysql.NewConnection(sql.DBConfig{
		Host:          cfg.MySQLHost,
		Port:          cfg.MySQLPort,
		Database:      cfg.MySQLDatabase,
		Username:      cfg.MySQLUsername,
		Password:      cfg.MySQLPassword,
		Params:        cfg.MySQLParams,
		MaxIdleConns:  cfg.MySQLMaxIdleConns,
		MaxOpenConns:  cfg.MySQLMaxOpenConns,
		SingularTable: cfg.MySQLSingularTable,
		DebugMode:     cfg.MySQLDebugMode,
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

	redisconn := redis.NewConnection(
		cfg.RedisHost,
		cfg.RedisPort,
		cfg.RedisUsername,
		cfg.RedisPassword,
		cfg.RedisNamespace,
		cfg.RedisDBNumber,
		cfg.RedisLocalCacheSize,
		cfg.RedisLocalCacheTTL,
		cfg.RedisMustAvailable,
		cfg.RedisDebugMode,
	)
	fmt.Println(redisconn)

	httpServer := http.NewServer(
		cfg.HTTPServerPort,
		cfg.HTTPServerEnableCORS,
		cfg.HTTPServerEnablePredefinedRoutes,
	)
	httpServer.CORS.AllowOrigins = cfg.HTTPServerAllowOrigins
	httpServer.CORS.AllowMethods = cfg.HTTPServerAllowMethods
	httpServer.CORS.AllowHeaders = cfg.HTTPServerAllowHeaders
	httpServer.CORS.MaxAge = cfg.HTTPServerMaxAge

	providerRepository := provider.NewRepository(dbconn, true)
	providerService := provider.NewService(providerRepository)
	providerHTTPHandler := api.NewProviderHTTPHandler(providerService)

	pokeapiClient := pokeapi.NewClient(cfg.PokeAPIAddressV2, cfg.PokeAPITimeout)
	pokemonService := pokemon.NewService(pokeapiClient)
	pokemonHTTPHandler := api.NewPokemonHTTPHandler(pokemonService)

	v1 := httpServer.Group("/v1")

	// Provider APIs:
	v1.POST("/provider", true, providerHTTPHandler.CreateProvider)
	v1.PUT("/provider/:uuid", true, providerHTTPHandler.UpdateProvider)
	v1.DELETE("/provider/:uuid", false, providerHTTPHandler.DeleteProviderByUUID)
	v1.GET("/provider/:uuid", false, providerHTTPHandler.GetProviderByUUID)
	v1.GET("/providers", false, providerHTTPHandler.GetProviders)

	// Pokemon APIs:
	v1.GET("/pokemon/:name", false, pokemonHTTPHandler.GetPokemonByName)

	promServer := prometheus.NewServer(
		cfg.PrometheusServerPort,
		cfg.PrometheusServerMetricsPath,
	)

	promServer.Monitor(
		&prometheus.Target{
			HTTPServer:    httpServer,
			ExcludePaths:  cfg.HTTPServerMonitorSkipPaths,
			GroupedStatus: cfg.HTTPServerMonitorGroupedStatus,
		},
	)

	server.RunServersGracefully(cfg.GracefulTimeout, promServer, httpServer)
}
