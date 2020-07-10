# go-gin-skeleton

A skeleton project for building RESTful API with Go &amp; Gin using Clean Architecture.

## Dependencies

* Web Framework: [gin-gonic/gin](https://github.com/gin-gonic/gin)
* REST Client: [go-resty/resty](https://github.com/go-resty/resty)
* Configuration: [kelseyhightower/envconfig](https://github.com/kelseyhightower/envconfig)
* Test Framework: [stretchr/testify](https://github.com/stretchr/testify)
* SQL Database ORM: [jinzhu/gorm](https://github.com/jinzhu/gorm)
* KV Store Abstrction: [philippgille/gokv](https://github.com/philippgille/gokv)
* Logging: [rs/zerolog](https://github.com/rs/zerolog)
* Prometheus: [prometheus/client_golang](https://github.com/prometheus/client_golang)
* Go HTTP Metrics: [slok/go-http-metrics](https://github.com/slok/go-http-metrics)
* Validator: [go-playground/validator](github.com/go-playground/validator)

## Docker

Use [docker-makefile](https://github.com/mvanholsteijn/docker-makefile) to build docker image with semantic versioning as its tag.
