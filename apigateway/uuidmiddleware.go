package apigateway

import (
	"context"
	"github.com/julienschmidt/httprouter"
	"github.com/satori/go.uuid"
	"net/http"
)

//httprouter middleware to generate and inject correlationId in the context
//this will be used in register fuc in http_server.go
func UUIDMiddleware(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		u2, err := uuid.NewV4()
		if err != nil {
			respondError(w, http.StatusInternalServerError, err)
		}
		ctx := r.Context()
		ctx = context.WithValue(ctx, "correlationid", u2)
		r = r.WithContext(ctx)
		next(w, r, ps)
	}
}
