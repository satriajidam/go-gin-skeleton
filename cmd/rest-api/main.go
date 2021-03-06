package main

import (
	"github.com/satriajidam/go-gin-skeleton/internal/config"
	"github.com/satriajidam/go-gin-skeleton/internal/service/api"
	"github.com/satriajidam/go-gin-skeleton/internal/service/client/pokeapi"
	"github.com/satriajidam/go-gin-skeleton/internal/service/pokemon"
	"github.com/satriajidam/go-gin-skeleton/internal/service/provider"
	"github.com/satriajidam/go-gin-skeleton/pkg/cache/redis"
	"github.com/satriajidam/go-gin-skeleton/pkg/database/sql"
	"github.com/satriajidam/go-gin-skeleton/pkg/database/sql/mysql"
	"github.com/satriajidam/go-gin-skeleton/pkg/server"
	"github.com/satriajidam/go-gin-skeleton/pkg/server/http"
	"github.com/satriajidam/go-gin-skeleton/pkg/server/prometheus"
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

	redisconn, err := redis.NewConnection(redis.RedisConfig{
		Host:      cfg.RedisHost,
		Port:      cfg.RedisPort,
		Username:  cfg.RedisUsername,
		Password:  cfg.RedisPassword,
		Namespace: cfg.RedisNamespace,
		DBNumber:  cfg.RedisDBNumber,
		DebugMode: cfg.RedisDebugMode,
	})
	defer func() {
		err := redisconn.Close()
		if err != nil {
			panic(err)
		}
	}()
	if err != nil && cfg.RedisMustAvailable {
		panic(err)
	}

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
	providerCache := provider.NewCache(redisconn)
	providerService := provider.NewService(providerRepository, providerCache)
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
