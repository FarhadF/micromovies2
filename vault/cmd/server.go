package main

import (
	"context"
	"fmt"
	"github.com/farhadf/micromovies2/vault"
	"github.com/farhadf/micromovies2/vault/pb"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	"github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	"github.com/julienschmidt/httprouter"
	"github.com/opentracing/opentracing-go"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	flag "github.com/spf13/pflag"
	"github.com/uber/jaeger-client-go/config"
	"google.golang.org/grpc"
	"io"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	jaegerZap "github.com/uber/jaeger-client-go/log/zap"
	"go.uber.org/zap"
)

func main() {
	var (
		httpAddr string
		gRPCAddr string
		console  bool
	)
	flag.StringVarP(&httpAddr, "http", "H", ":8086", "http listen address")
	flag.StringVarP(&gRPCAddr, "grpc", "G", ":8085", "gRPC listen address")
	flag.BoolVarP(&console, "console", "c", false, "turns on pretty console logging")
	flag.Parse()
	ctx := context.Background()
	//zap
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	fieldKeys := []string{"method", "error"}
	requestCount := kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
		Namespace: "micromovies2",
		Subsystem: "vault_service",
		Name:      "request_count",
		Help:      "Number of requests received.",
	}, fieldKeys)
	requestLatency := kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
		Namespace: "micromovies2",
		Subsystem: "vault_service",
		Name:      "request_latency_seconds",
		Help:      "Total duration of requests in seconds.",
	}, fieldKeys)

	svc := vault.NewService()
	svc = vault.LoggingMiddleware{logger, svc}
	svc = vault.InstrumentingMiddleware{requestCount, requestLatency, svc}
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
		HashEndpoint:     hashEndpoint,
		ValidateEndpoint: validateEndpoint,
	}
	//tracing
	tracer, closer := initJaeger("vaultService", logger)
	defer closer.Close()
	opentracing.SetGlobalTracer(tracer)
	span := tracer.StartSpan("server")
	defer span.Finish()
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
		logger.Info("", zap.String("grpc:", gRPCAddr))
		handler := vault.NewGRPCServer(ctx, endpoints)
		//add grpc_opentracing interceptor for server
		gRPCServer := grpc.NewServer(grpc.UnaryInterceptor(grpc_opentracing.UnaryServerInterceptor()))
		pb.RegisterVaultServer(gRPCServer, handler)
		errChan <- gRPCServer.Serve(listener)
	}()

	logger.Info(gRPCAddr)
	logger.Fatal("", zap.Error(<-errChan))
}

// initJaeger returns an instance of Jaeger Tracer that samples 100% of traces and logs all spans to stdout.
func initJaeger(service string, logger *zap.Logger) (opentracing.Tracer, io.Closer) {
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
	tracer, closer, err := cfg.NewTracer(config.Logger(jaegerZap.NewLogger(logger)))
	if err != nil {
		logger.Panic("ERROR: cannot init Jaeger", zap.Error(err))
	}
	return tracer, closer
}
