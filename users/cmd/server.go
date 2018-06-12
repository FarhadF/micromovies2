package main

import (
	"context"
	"fmt"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	"github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	"github.com/jackc/pgx"
	"github.com/julienschmidt/httprouter"
	"github.com/opentracing/opentracing-go"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	flag "github.com/spf13/pflag"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"io"
	"micromovies2/users"
	"micromovies2/users/pb"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

var pool *pgx.ConnPool

func main() {
	//zap
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	var (
		console  bool
		httpAddr string
		gRPCAddr string
	)
	flag.StringVarP(&httpAddr, "http", "H", ":8083", "http listen address")
	flag.StringVarP(&gRPCAddr, "grpc", "g", ":8084", "GRPC Address")
	flag.BoolVarP(&console, "console", "c", false, "turns on pretty console logging")
	flag.Parse()
	logger.Info("starting grpc server at" + gRPCAddr)
	ctx := context.Background()
	//database
	connPoolConfig := pgx.ConnPoolConfig{
		ConnConfig: pgx.ConnConfig{
			Host:     "127.0.0.1",
			Port:     26257,
			User:     "app_user",
			Database: "app_database",
			//Logger: logger, todo: fix logger
		},
		MaxConnections: 5,
	}
	pool, err := pgx.NewConnPool(connPoolConfig)
	if err != nil {
		logger.Fatal("Unable to create connection pool", zap.Error(err))
	}
	/*db, err := sql.Open("postgres", "postgresql://app_user@localhost:26257/app_database?sslmode=disable")
	if err != nil {
		logger.Fatal().Err(err).Msg("db connection failed")
	}*/

	//instrumentation
	fieldKeys := []string{"method", "error"}
	requestCount := kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
		Namespace: "my_group",
		Subsystem: "users_service",
		Name:      "request_count",
		Help:      "Number of requests received.",
	}, fieldKeys)
	requestLatency := kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
		Namespace: "my_group",
		Subsystem: "users_service",
		Name:      "request_latency_microseconds",
		Help:      "Total duration of requests in microseconds.",
	}, fieldKeys)

	// init users service
	var svc users.Service
	svc = users.NewService(pool, *logger)
	//wire logging
	svc = users.LoggingMiddleware{*logger, svc}
	//wire instrumentation
	svc = users.InstrumentingMiddleware{requestCount, requestLatency, svc}
	errChan := make(chan error)
	// creating Endpoints struct
	endpoints := users.Endpoints{
		NewUserEndpoint:        users.MakeNewUserEndpoint(svc),
		GetUserByEmailEndpoint: users.MakeGetUserByEmailEndpoint(svc),
		ChangePasswordEndpoint: users.MakeChangePasswordEndpoint(svc),
		LoginEndpoint:          users.MakeLoginEndpoint(svc),
	}
	//tracing
	tracer, closer := initJaeger("userService")
	defer closer.Close()
	opentracing.SetGlobalTracer(tracer)
	//span := tracer.StartSpan("server")
	//defer span.Finish()
	//execute grpc server
	go func() {
		listener, err := net.Listen("tcp", gRPCAddr)
		if err != nil {
			errChan <- err
			return
		}
		handler := users.NewGRPCServer(ctx, endpoints)
		//grpc server with grpc-ecosystem/go-grpc-middleware interceptor
		grpcServer := grpc.NewServer(grpc.UnaryInterceptor(grpc_opentracing.UnaryServerInterceptor()))
		pb.RegisterUsersServer(grpcServer, handler)
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
	logger.Error("", zap.Error(<-errChan))
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
