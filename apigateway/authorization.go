package apigateway

import (
	"context"
	"errors"
	"fmt"
	"github.com/casbin/casbin"
	"github.com/julienschmidt/httprouter"
	"micromovies2/jwtauth"
	"micromovies2/jwtauth/client"
	"net/http"
)

// Authorizer is a middleware for authorization
func (a *Authorizer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//todo: exclude public urls
	// check exclude url
	if len(a.excludeUrl) > 0 {
		for _, url := range a.excludeUrl {
			fmt.Println(url, r.URL.Path)
			if url == r.URL.Path {
				a.next.ServeHTTP(w, r)
				return
			}
		}
	}
	// check exclude url prefix
	/*	if len(config.ExcludePrefix) > 0 {
		for _, prefix := range config.ExcludePrefix {
			if strings.HasPrefix(c.Req.URL.Path, prefix) {
				c.Next()
				return
			}
		}
	}*/
	auth := &Authorizer{enforcer: a.enforcer}
	//extract token
	token := auth.getToken(r)
	var role string
	if token == "" {
		role = "public"
	} else {
		//parse and validate token
		claims, err := auth.getClaims(a.ctx, a.jwtAuthService, token)
		if err != nil {
			respondError(w, http.StatusForbidden, err)
			return
		}
		role = claims.Role
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

// Authorizer stores the casbin handler
type Authorizer struct {
	ctx            context.Context
	next           *httprouter.Router
	enforcer       *casbin.Enforcer
	jwtAuthService jwtauth.Service
	excludeUrl     []string
}

// Make a constructor for our middleware type since its fields are not exported (in lowercase)
func NewAuthMiddleware(ctx context.Context, next *httprouter.Router, e *casbin.Enforcer, jwtAuthService jwtauth.Service, excludeUrls []string) *Authorizer {
	return &Authorizer{ctx: ctx, next: next, enforcer: e, jwtAuthService: jwtAuthService, excludeUrl: excludeUrls}
}

// GetToken gets the jwt token from the request.
func (a *Authorizer) getToken(r *http.Request) string {
	token := r.Header.Get("Authorization")
	return token
}

// GetClaims gets the  role from jwt claims.
func (a *Authorizer) getClaims(ctx context.Context, jwtAuthService jwtauth.Service, token string) (jwtauth.Claims, error) {
	claims, err := client.ParseToken(ctx, jwtAuthService, token)
	return claims, err
}
