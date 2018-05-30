package main

import (
	"context"
	"github.com/casbin/casbin"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	"github.com/julienschmidt/httprouter"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	flag "github.com/spf13/pflag"
	"go.uber.org/zap"
	"micromovies2/apigateway"
	"micromovies2/jwtauth"
	"net/http"
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
	//zap
	logger, _ := zap.NewProduction()
	defer logger.Sync()
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
	svc = apigateway.LoggingMiddleware{*logger, svc}
	svc = apigateway.InstrumentingMiddleware{requestCount, requestLatency, svc}

	// setup casbin auth rules
	e, err := casbin.NewEnforcerSafe("/home/balrog/go/src/micromovies2/apigateway/cmd/model.conf", "/home/balrog/go/src/micromovies2/apigateway/cmd/policy.csv", false)
	//disable casbin log
	e.EnableLog(false)
	if err != nil {
		logger.Fatal("", zap.Error(err))
	}
	jwtAuthService := jwtauth.NewService()
	// HTTP transport

	logger.Info("", zap.String("http:", httpAddr))
	//httprouter
	r := httprouter.New()
	apigateway.Endpoints{
		Ctx:              ctx,
		LoginEndpoint:    apigateway.MakeLoginEndpoint(svc),
		RegisterEndpoint: apigateway.MakeRegisterEndpoint(svc),
	}.Register(r)
	authMiddleware := apigateway.NewAuthMiddleware(ctx, r, e, jwtAuthService)
	logger.Fatal("", zap.Error(http.ListenAndServe(httpAddr, authMiddleware)))
}
