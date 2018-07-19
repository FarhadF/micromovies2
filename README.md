Micromovies2
=======================================================================

Micromovies2 is a sample application using [Go-Kit](https://github.com/go-kit/kit) and various dependencies, focusing on microservices architecture as well as deployment, monitoring, tracing, logging and stress testing. [CockroachDB](www.cockroachlabs.com/‎) is used as RDBMS. It's all Go! ecosystem aside from [jmeter](http://jmeter.apache.org/) for stress testing.

- [Microservices](#Microservices)
- [Notable Packages and Systems](#notable-packages-and-systems)
- [Quickstart](#quickstart)

## Microservices

Micromovies2 is divided into 5 microservices, each having it's purpose. Microservices internal communication is provided by grpc. Client facing communication is using REST API and is provided by APIGateway Service. Here are the services:

- APIGateway
- JWTAuth
- Movies
- Users
- Vault

## Notable Packages and Systems

- go-kit
- httprouter
- zap
- opentracing/opentracing-go
- swagger
- casbin
- pflags
- grpc
- cockroachdb
- jaeger
- prometheus
- grafana
- jmeter
- vgo


## Quickstart

1. ```git clone https://github.com/farhadf/micromovies2```
2. ```cd micromovies2```
2. ```docker-compose up```


