package apigateway

import (
	"context"
	"errors"
	"github.com/casbin/casbin"
	"github.com/farhadf/micromovies2/jwtauth"
	"github.com/farhadf/micromovies2/jwtauth/client"
	"github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	"github.com/julienschmidt/httprouter"
	"google.golang.org/grpc"
	"net/http"
	"strings"
)

//overload the ServeHTTP method
func (a *Authorizer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// check exclude url
	if len(a.excludeUrl) > 0 {
		for _, url := range a.excludeUrl {
			if url == r.URL.Path {
				r = r.WithContext(a.ctx)
				a.next.ServeHTTP(w, r)
				return
			}
		}
	}
	// check exclude url prefix
	if len(a.excludePrefix) > 0 {
		for _, prefix := range a.excludePrefix {
			if strings.HasPrefix(r.URL.Path, prefix) {
				r = r.WithContext(a.ctx)
				a.next.ServeHTTP(w, r)
				return
			}
		}
	}
	//auth := &Authorizer{enforcer: a.enforcer}
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
		claims, err := a.getClaims(a.ctx, token[1])
		if err != nil {
			respondError(w, http.StatusForbidden, err)
			return
		}
		role = claims.Role
		//put desired data in the context
		a.ctx = context.WithValue(a.ctx, "email", claims.Email)
		a.ctx = context.WithValue(a.ctx, "role", claims.Role)
		//put the context in the http.request context and make sure you take it out at http handlers
		r = r.WithContext(a.ctx)
	}
	// casbin enforce
	res, err := a.enforcer.EnforceSafe(role, r.URL.Path, r.Method)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	if res {
		a.next.ServeHTTP(w, r)
	} else {
		respondError(w, http.StatusForbidden, errors.New("unauthorized"))
		return
	}
}

// Authorizer is a middleware for authorization
// Authorizer stores the casbin handler plus everything we need to feed to ServeHTTP
type Authorizer struct {
	ctx           context.Context
	next          *httprouter.Router
	enforcer      *casbin.Enforcer
	excludeUrl    []string
	excludePrefix []string
	jwtAuthAddr   string
}

// Make a constructor for our middleware type since its fields are not exported
func NewAuthMiddleware(ctx context.Context, next *httprouter.Router, e *casbin.Enforcer,
	excludeUrls []string, jwtAuthAddr string) *Authorizer {
	return &Authorizer{ctx: ctx, next: next, enforcer: e, excludeUrl: excludeUrls, jwtAuthAddr: jwtAuthAddr}
}

// GetToken gets the jwt token from the request.
func (a *Authorizer) getToken(r *http.Request) string {
	token := r.Header.Get("Authorization")
	return token
}

//todo: fix grpc call to jwt service
// GetClaims gets the  role from jwt claims.
func (a *Authorizer) getClaims(ctx context.Context, token string) (jwtauth.Claims, error) {
	conn, err := grpc.Dial(a.jwtAuthAddr, grpc.WithInsecure(), grpc.WithUnaryInterceptor(grpc_opentracing.UnaryClientInterceptor()))
	if err != nil {
		return jwtauth.Claims{}, err
	}
	defer conn.Close()
	jwtAuthService := client.NewGRPCClient(conn)
	claims, err := client.ParseToken(ctx, jwtAuthService, token)
	return claims, err
}
