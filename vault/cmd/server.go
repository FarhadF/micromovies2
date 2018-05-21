package main

import (
	flag "github.com/spf13/pflag"
	"micromovies2/vault"
	"context"
	"os"
	"os/signal"
	"syscall"
	"fmt"
	"net/http"
	"github.com/rs/zerolog"
	"net"
	"micromovies2/vault/pb"
	"google.golang.org/grpc"
	"github.com/rs/zerolog/log"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	"github.com/julienschmidt/httprouter"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	var (
		httpAddr string
		gRPCAddr string
		console  bool
	)
	flag.StringVarP(&httpAddr, "http", "H", ":8082", "http listen address")
	flag.StringVarP(&gRPCAddr, "grpc", "G", ":8081", "gRPC listen address")
	flag.BoolVarP(&console, "console", "c", false, "turns on pretty console logging")
	flag.Parse()
	ctx := context.Background()
	//zerolog
	logger := zerolog.New(os.Stderr).With().Timestamp().Caller().Logger()
	//console pretty printing
	if console {
		logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}
	fieldKeys := []string{"method", "error"}
	requestCount := kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
		Namespace: "my_group",
		Subsystem: "vault_service",
		Name:      "request_count",
		Help:      "Number of requests received.",
	}, fieldKeys)
	requestLatency := kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
		Namespace: "my_group",
		Subsystem: "vault_service",
		Name:      "request_latency_microseconds",
		Help:      "Total duration of requests in microseconds.",
	}, fieldKeys)
	countResult := kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
		Namespace: "my_group",
		Subsystem: "vault_service",
		Name:      "count_result",
		Help:      "The result of each count method.",
	}, []string{}) // no fields here

	svc := vault.NewService()
	svc = vault.LoggingMiddleware{logger, svc}
	svc = vault.InstrumentingMiddleware{requestCount, requestLatency, countResult, svc}
	errChan := make(chan error)
	//os signal handling
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errChan <- fmt.Errorf("%s", <-c)
	}()
	hashEndpoint := vault.MakeHashEndpoint(svc)
	validateEndpoint := vault.MakeValidateEndpoint(svc)
	endpoints := vault.Endpoints{
		HashEndpoint:
		hashEndpoint,
		ValidateEndpoint: validateEndpoint,
	}
	// HTTP transport
	go func() {
		//httprouter initialization
		router := httprouter.New()
		//handler will be used for net/http handle compatibility
		router.Handler("GET", "/metrics", promhttp.Handler())
		errChan <- http.ListenAndServe(httpAddr, router)
	}()
	// GRPC transport
	go func() {
		listener, err := net.Listen("tcp", gRPCAddr)
		if err != nil {
			errChan <- err
			return
		}
		logger.Info().Str("grpc:", gRPCAddr).Msg("")
		handler := vault.NewGRPCServer(ctx, endpoints)
		gRPCServer := grpc.NewServer()
		pb.RegisterVaultServer(gRPCServer, handler)
		errChan <- gRPCServer.Serve(listener)
	}()

	logger.Info().Msg(gRPCAddr)
	logger.Fatal().Err(<-errChan).Msg("")
}
