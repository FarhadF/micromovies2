package main

import (
	"context"
	"fmt"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	"github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	"github.com/julienschmidt/httprouter"
	"github.com/opentracing/opentracing-go"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	flag "github.com/spf13/pflag"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
	"google.golang.org/grpc"
	"io"
	"micromovies2/jwtauth"
	"micromovies2/jwtauth/pb"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
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
	svc = jwtauth.InstrumentingMiddleware{requestCount, requestLatency, svc}
	errChan := make(chan error)
	// creating Endpoints struct
	endpoints := jwtauth.Endpoints{
		GenerateTokenEndpoint: jwtauth.MakeGenerateTokenEndpoint(svc),
		ParseTokenEndpoint:    jwtauth.MakeParseTokenEndpoint(svc),
	}
	//tracing
	tracer, closer := initJaeger("jwtService")
	defer closer.Close()
	opentracing.SetGlobalTracer(tracer)
	span := tracer.StartSpan("server")
	defer span.Finish()
	//execute grpc server
	go func() {
		listener, err := net.Listen("tcp", gRPCAddr)
		if err != nil {
			errChan <- err
			return
		}
		handler := jwtauth.NewGRPCServer(ctx, endpoints)
		//add grpc_opentracing interceptor for server
		grpcServer := grpc.NewServer(grpc.UnaryInterceptor(grpc_opentracing.UnaryServerInterceptor()))
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

// initJaeger returns an instance of Jaeger Tracer that samples 100% of traces and logs all spans to stdout.
func initJaeger(service string) (opentracing.Tracer, io.Closer) {
	cfg := &config.Configuration{
		Sampler: &config.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &config.ReporterConfig{
			LogSpans: false,
		},
		ServiceName: service,
	}
	tracer, closer, err := cfg.NewTracer(config.Logger(jaeger.StdLogger))
	if err != nil {
		panic(fmt.Sprintf("ERROR: cannot init Jaeger: %v\n", err))
	}
	return tracer, closer
}
