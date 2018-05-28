package main

import (
	flag "github.com/spf13/pflag"
	"micromovies2/apigateway"
	"context"
	"os"
	"os/signal"
	"syscall"
	"fmt"
	"net/http"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	"github.com/julienschmidt/httprouter"
)

func main() {
	var (
		httpAddr string
		console  bool
	)
	flag.StringVarP(&httpAddr, "http", "H", ":8089", "http listen address")
	flag.BoolVarP(&console, "console", "c", false, "turns on pretty console logging")
	flag.Parse()
	ctx := context.Background()
	//zerolog
	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()
	//console pretty printing
	if console {
		logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}
	fieldKeys := []string{"method", "error"}
	requestCount := kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
		Namespace: "my_group",
		Subsystem: "apigateway_service",
		Name:      "request_count",
		Help:      "Number of requests received.",
	}, fieldKeys)
	requestLatency := kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
		Namespace: "my_group",
		Subsystem: "apigateway_service",
		Name:      "request_latency_microseconds",
		Help:      "Total duration of requests in microseconds.",
	}, fieldKeys)

	svc := apigateway.NewService()
	svc = apigateway.LoggingMiddleware{logger, svc}
	svc = apigateway.InstrumentingMiddleware{requestCount, requestLatency,  svc}
	errChan := make(chan error)
	//os signal handling
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errChan <- fmt.Errorf("%s", <-c)
	}()
	// HTTP transport
	go func() {
		logger.Info().Str("http:", httpAddr).Msg("")
		//handler := apigateway.NewHTTPServer(ctx, endpoints)
		//httprouter
		r := httprouter.New()
		apigateway.Endpoints{
			Ctx: ctx,
			// This is incredibly laborious when we want to add e.g. rate
			// limiters. It would be better to bundle all the endpoints up,
			// somehow... or, use code generation, of course.
			LoginEndpoint:     apigateway.MakeLoginEndpoint(svc),
		}.Register(r)
		errChan <- http.ListenAndServe(httpAddr, r)
	}()
	logger.Fatal().Err(<-errChan).Msg("")
}