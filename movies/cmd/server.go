package main

import (
	"context"
	"fmt"
	"github.com/farhadf/micromovies2/movies"
	"github.com/farhadf/micromovies2/movies/pb"
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
		console  bool
		httpAddr string
		gRPCAddr string
		dbHost   string
		dbPort   uint16
		dbName   string
		dbUser   string
	)
	flag.StringVarP(&httpAddr, "http", "H", ":8082", "http listen address")
	flag.StringVarP(&gRPCAddr, "grpc", "g", ":8081", "GRPC Address")
	flag.StringVarP(&dbHost, "dbhost", "d", "127.0.0.1", "Database Hostname/IP Address")
	flag.Uint16VarP(&dbPort, "dbport", "p", 26257, "Database port number")
	flag.StringVarP(&dbName, "dbname", "n", "app_database", "Database name")
	flag.StringVarP(&dbUser, "dbuser", "u", "app_user", "Database user")
	flag.BoolVarP(&console, "console", "c", false, "turns on pretty console logging")
	flag.Parse()
	logger.Info("starting grpc server at" + string(gRPCAddr))
	ctx := context.Background()
	//database
	connPoolConfig := pgx.ConnPoolConfig{
		ConnConfig: pgx.ConnConfig{
			Host:     dbHost,
			Port:     dbPort,
			User:     dbUser,
			Database: dbName,
			Logger:   zapadapter.NewLogger(logger),
			LogLevel: pgx.LogLevelWarn,
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
		Subsystem: "movies_service",
		Name:      "request_count",
		Help:      "Number of requests received.",
	}, fieldKeys)
	requestLatency := kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
		Namespace: "my_group",
		Subsystem: "movies_service",
		Name:      "request_latency_microseconds",
		Help:      "Total duration of requests in microseconds.",
	}, fieldKeys)

	// init movies service
	var svc movies.Service
	svc = movies.NewService(pool, logger)
	//wire logging
	svc = movies.LoggingMiddleware{logger, svc}
	//wire instrumentation
	svc = movies.InstrumentingMiddleware{requestCount, requestLatency, svc}
	errChan := make(chan error)
	// creating Endpoints struct
	endpoints := movies.Endpoints{
		GetMoviesEndpoint:    movies.MakeGetMoviesEndpoint(svc),
		GetMovieByIdEndpoint: movies.MakeGetMovieByIdEndpoint(svc),
		NewMovieEndpoint:     movies.MakeNewMovieEndpoint(svc),
		DeleteMovieEndpoint:  movies.MakeDeleteMovieEndpoint(svc),
		UpdateMovieEndpoint:  movies.MakeUpdateMovieEndpoint(svc),
	}
	//tracing
	tracer, closer := initJaeger("moviesService", logger)
	defer closer.Close()
	opentracing.SetGlobalTracer(tracer)
	//execute grpc server
	go func() {
		listener, err := net.Listen("tcp", gRPCAddr)
		if err != nil {
			errChan <- err
			return
		}
		handler := movies.NewGRPCServer(ctx, endpoints)
		grpcServer := grpc.NewServer(grpc.UnaryInterceptor(grpc_opentracing.UnaryServerInterceptor()))
		pb.RegisterMoviesServer(grpcServer, handler)
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
