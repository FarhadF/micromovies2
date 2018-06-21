// APIGateway
//
// This documentation describes APIGateway APIs
//
//     Schemes: http
//     BasePath: /v1
//     Version: 1.0.0
//     License: MIT http://opensource.org/licenses/MIT
//     Contact: Farhad Farahi<ff@ff.f>
//     Host: localhost
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Security:
//     - bearer
//
//     SecurityDefinitions:
//     bearer:
//          type: JWT
//          name: Authorization
//          in: header
//
// swagger:meta
package main

import (
	"context"
	"fmt"
	"github.com/casbin/casbin"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	"github.com/julienschmidt/httprouter"
	"github.com/opentracing/opentracing-go"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	flag "github.com/spf13/pflag"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
	"go.uber.org/zap"
	"io"
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
	//tracing
	tracer, closer := initJaeger("api-gateway")
	defer closer.Close()
	opentracing.SetGlobalTracer(tracer)
	// HTTP transport
	logger.Info("", zap.String("http:", httpAddr))
	//httprouter
	r := httprouter.New()
	apigateway.Endpoints{
		Ctx:                    ctx,
		LoginEndpoint:          apigateway.MakeLoginEndpoint(svc),
		RegisterEndpoint:       apigateway.MakeRegisterEndpoint(svc),
		ChangePasswordEndpoint: apigateway.MakeChangePasswordEndpoint(svc),
		GetMovieByIdEndpoint:   apigateway.MakeGetMovieByIdEndpoint(svc),
	}.Register(r)
	excludeUrls := []string{"/v1/login", "/v1/register"}
	authMiddleware := apigateway.NewAuthMiddleware(ctx, r, e, jwtAuthService, excludeUrls)
	logger.Fatal("", zap.Error(http.ListenAndServe(httpAddr, authMiddleware)))
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
