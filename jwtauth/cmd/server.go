package main

import (
	"os"
	flag "github.com/spf13/pflag"
	"net"
	"github.com/julienschmidt/httprouter"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"os/signal"
	"syscall"
	"fmt"
	"context"
	"github.com/rs/zerolog"
	"zerolog/log"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	"micromovies2/jwtauth"
	"google.golang.org/grpc"
	"micromovies2/jwtauth/pb"
)

func main() {
	//zerolog
	logger := zerolog.New(os.Stderr).With().Timestamp().Caller().Logger()
	var (
		console  bool
		httpAddr string
		gRPCAddr string
	)
	flag.StringVarP(&httpAddr, "http", "H", ":8088", "http listen address")
	flag.StringVarP(&gRPCAddr, "grpc", "g", ":8087", "GRPC Address")
	flag.BoolVarP(&console, "console", "c", false, "turns on pretty console logging")
	flag.Parse()
	logger.Info().Msg("starting grpc server at" + string(gRPCAddr))
	if console {
		logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}
	ctx := context.Background()
	//instrumentation
	fieldKeys := []string{"method", "error"}
	requestCount := kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
		Namespace: "my_group",
		Subsystem: "jwt_service",
		Name:      "request_count",
		Help:      "Number of requests received.",
	}, fieldKeys)
	requestLatency := kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
		Namespace: "my_group",
		Subsystem: "jwt_service",
		Name:      "request_latency_microseconds",
		Help:      "Total duration of requests in microseconds.",
	}, fieldKeys)

	// init jwt service
	svc := jwtauth.NewService()
	//wire logging
	svc = jwtauth.LoggingMiddleware{logger, svc}
	//wire instrumentation
	svc = jwtauth.InstrumentingMiddleware{requestCount, requestLatency,  svc}
	errChan := make(chan error)
	// creating Endpoints struct
	endpoints := jwtauth.Endpoints{
		GenerateTokenEndpoint: jwtauth.MakeGenerateTokenEndpoint(svc),
	}
	//execute grpc server
	go func() {
		listener, err := net.Listen("tcp", gRPCAddr)
		if err != nil {
			errChan <- err
			return
		}
		handler := jwtauth.NewGRPCServer(ctx, endpoints)
		grpcServer := grpc.NewServer()
		pb.RegisterJWTServer(grpcServer, handler)
		errChan <- grpcServer.Serve(listener)
	}()
	// HTTP transport
	go func() {
		//httprouter initialization
		router := httprouter.New()
		//handler will be used for net/http handle compatibility
		router.Handler("GET", "/metrics", promhttp.Handler())
		errChan <- http.ListenAndServe(httpAddr, router)
	}()
	//Handle os signals
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errChan <- fmt.Errorf("%s", <-c)
	}()
	logger.Error().Err(<-errChan).Msg("")
}

