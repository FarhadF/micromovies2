package main

import (
	"context"
	"fmt"
	"github.com/farhadf/micromovies2/users"
	"github.com/farhadf/micromovies2/users/pb"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	"github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/log/zapadapter"
	"github.com/julienschmidt/httprouter"
	"github.com/opentracing/opentracing-go"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	flag "github.com/spf13/pflag"
	"github.com/uber/jaeger-client-go/config"
	jaegerZap "github.com/uber/jaeger-client-go/log/zap"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"io"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	//zap
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	var (
		console     bool
		httpAddr    string
		gRPCAddr    string
		dbHost      string
		dbPort      uint16
		dbName      string
		dbUser      string
		vaultAddr   string
		jwtAuthAddr string
	)
	flag.StringVarP(&httpAddr, "http", "H", ":8083", "http listen address")
	flag.StringVarP(&gRPCAddr, "grpc", "g", ":8084", "GRPC Address")
	flag.StringVarP(&dbHost, "dbhost", "d", "127.0.0.1", "Database Hostname/IP Address")
	flag.Uint16VarP(&dbPort, "dbport", "p", 26257, "Database port number")
	flag.StringVarP(&dbName, "dbname", "n", "app_database", "Database name")
	flag.StringVarP(&dbUser, "dbuser", "u", "app_user", "Database user")
	flag.StringVarP(&vaultAddr, "vault", "v", ":8085", "vault service listen address, example: localhost:8085")
	flag.StringVarP(&jwtAuthAddr, "jwtauth", "j", ":8087", "jwtAuth service listen address, example: localhost:8087")
	flag.BoolVarP(&console, "console", "c", false, "turns on pretty console logging")
	flag.Parse()
	logger.Info("starting grpc server at" + gRPCAddr)
	ctx := context.Background()
	//database
	connPoolConfig := pgx.ConnPoolConfig{
		ConnConfig: pgx.ConnConfig{
			Host:     dbHost,
			Port:     dbPort,
			User:     dbUser,
			Database: dbName,
			Logger:   zapadapter.NewLogger(logger),
		},
		MaxConnections: 5,
	}
	pool, err := pgx.NewConnPool(connPoolConfig)
	if err != nil {
		logger.Fatal("Unable to create connection pool", zap.Error(err))
	}
	config := users.Config{
		VaultAddr:   vaultAddr,
		JwtAuthAddr: jwtAuthAddr,
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
	svc = users.NewService(pool, *logger, config)
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
	tracer, closer := initJaeger("userService", logger)
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
	//Type Logger is an adapter from zap Logger to jaeger-lib Logger. New logger will actually do this for us.
	tracer, closer, err := cfg.NewTracer(config.Logger(jaegerZap.NewLogger(logger)))
	if err != nil {
		logger.Panic("ERROR: cannot init Jaeger:", zap.Error(err))
	}
	return tracer, closer
}
