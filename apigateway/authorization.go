package apigateway

import (
	"context"
	"errors"
	"github.com/casbin/casbin"
	"github.com/farhadf/micromovies2/jwtauth"
	"github.com/farhadf/micromovies2/jwtauth/client"
	"github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	"github.com/julienschmidt/httprouter"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"
	"net/http"
	"strings"
)

// Authorizer is a middleware for authorization
// Authorizer stores the casbin handler plus everything we need to feed to Authorization Middleware
type Authorizer struct {
	ctx         context.Context
	enforcer    *casbin.Enforcer
	jwtAuthAddr string
}

// Make a constructor for our middleware type since its fields are not exported
func NewAuthMiddleware(ctx context.Context, e *casbin.Enforcer,
	jwtAuthAddr string) *Authorizer {
	return &Authorizer{ctx: ctx, enforcer: e, jwtAuthAddr: jwtAuthAddr}
}

// GetToken gets the jwt token from the request.
func (a *Authorizer) getToken(r *http.Request) string {
	token := r.Header.Get("Authorization")
	return token
}

// GetClaims gets the  role from jwt claims.
func (a *Authorizer) getClaims(ctx context.Context, token string) (jwtauth.Claims, error) {
	if span := opentracing.SpanFromContext(ctx); span != nil {
		span := span.Tracer().StartSpan("GetClaims", opentracing.ChildOf(span.Context()))
		defer span.Finish()
		ctx = opentracing.ContextWithSpan(ctx, span)
	}
	conn, err := grpc.Dial(a.jwtAuthAddr, grpc.WithInsecure(), grpc.WithUnaryInterceptor(grpc_opentracing.UnaryClientInterceptor()))
	if err != nil {
		return jwtauth.Claims{}, err
	}
	defer conn.Close()
	jwtAuthService := client.NewGRPCClient(conn)
	claims, err := client.ParseToken(ctx, jwtAuthService, token)
	return claims, err
}

func (a *Authorizer) AuthorizationMiddleware(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		//get request context
		ctx := r.Context()
		if span := opentracing.SpanFromContext(ctx); span != nil {
			span := span.Tracer().StartSpan("AuthorizationMiddleware", opentracing.ChildOf(span.Context()))
			defer span.Finish()
			ctx = opentracing.ContextWithSpan(ctx, span)
		}
		//extract token
		authHeader := a.getToken(r)
		var role string
		if authHeader == "" {
			role = "public"
		} else {
			token := strings.Split(authHeader, " ")
			if len(token) != 2 || token[0] != "Bearer" {
				respondError(w, http.StatusBadRequest, errors.New("bad token"))
				return
			}
			//parse and validate token
			claims, err := a.getClaims(ctx, token[1])
			if err != nil {
				respondError(w, http.StatusForbidden, err)
				return
			}
			role = claims.Role
			//put desired data in the context
			ctx = context.WithValue(ctx, "email", claims.Email)
			ctx = context.WithValue(ctx, "role", claims.Role)
			//put the context in the http.request context and make sure you take it out at http handlers
			r = r.WithContext(ctx)
		}
		// casbin enforce
		res, err := a.enforcer.EnforceSafe(role, r.URL.Path, r.Method)
		if err != nil {
			respondError(w, http.StatusInternalServerError, err)
			return
		}
		if res {
			r = r.WithContext(ctx)
			next(w, r, ps)
			return
		} else {
			respondError(w, http.StatusForbidden, errors.New("unauthorized"))
			return
		}
	}
}
