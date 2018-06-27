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
	"github.com/farhadf/micromovies2/apigateway"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	"github.com/julienschmidt/httprouter"
	"github.com/opentracing/opentracing-go"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	flag "github.com/spf13/pflag"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
	"go.uber.org/zap"
	"io"
	"net/http"
	"path/filepath"
)

func main() {
	var (
		httpAddr    string
		moviesAddr  string
		usersAddr   string
		jwtAuthAddr string
		console     bool
	)
	flag.StringVarP(&httpAddr, "http", "H", ":8089", "http listen address, example: localhost:8089")
	flag.StringVarP(&moviesAddr, "movies", "m", ":8081", "movies service listen address, example: localhost:8081")
	flag.StringVarP(&usersAddr, "users", "u", ":8084", "users service listen address, example: localhost:8084")
	flag.StringVarP(&jwtAuthAddr, "jwtauth", "j", ":8087", "jwtAuth service listen address, example: localhost:8087")
	flag.BoolVarP(&console, "console", "c", false, "turns on pretty console logging")
	flag.Parse()
	config := apigateway.Config{
		MoviesAddr:  moviesAddr,
		UsersAddr:   usersAddr,
		JwtAuthAddr: jwtAuthAddr,
	}
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

	svc := apigateway.NewService(config)
	svc = apigateway.LoggingMiddleware{*logger, svc}
	svc = apigateway.InstrumentingMiddleware{requestCount, requestLatency, svc}
	// setup casbin auth rules
	modelPath, _ := filepath.Abs("./model.conf")
	policyPath, _ := filepath.Abs("./policy.csv")
	e, err := casbin.NewEnforcerSafe(modelPath, policyPath, false)
	//disable casbin log
	e.EnableLog(false)
	if err != nil {
		logger.Fatal("", zap.Error(err))
	}
	//tracing
	tracer, closer := initJaeger("api-gateway")
	defer closer.Close()
	opentracing.SetGlobalTracer(tracer)
	// HTTP transport
	logger.Info("", zap.String("http:", httpAddr))
	//httprouter
	r := httprouter.New()
	//authorizer
	//excludeUrls := []string{"/v1/login", "/v1/register"}
	authMiddleware := apigateway.NewAuthMiddleware(ctx, e, jwtAuthAddr)
	apigateway.Endpoints{
		Ctx:                    ctx,
		LoginEndpoint:          apigateway.MakeLoginEndpoint(svc),
		RegisterEndpoint:       apigateway.MakeRegisterEndpoint(svc),
		ChangePasswordEndpoint: apigateway.MakeChangePasswordEndpoint(svc),
		GetMovieByIdEndpoint:   apigateway.MakeGetMovieByIdEndpoint(svc),
		NewMovieEndpoint:       apigateway.MakeNewMovieEndpoint(svc),
	}.Register(r, *authMiddleware)

	logger.Fatal("", zap.Error(http.ListenAndServe(httpAddr, r)))
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
		//todo: use structured logging
		panic(fmt.Sprintf("ERROR: cannot init Jaeger: %v\n", err))
	}
	return tracer, closer
}
