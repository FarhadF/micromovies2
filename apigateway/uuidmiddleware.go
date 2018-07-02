package apigateway

import (
	"context"
	"github.com/julienschmidt/httprouter"
	"github.com/opentracing/opentracing-go"
	"github.com/satori/go.uuid"
	"net/http"
)

//httprouter middleware to generate and inject correlationId in the context
//this will be used in register fuc in http_server.go
func UUIDMiddleware(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		//get request context
		ctx := r.Context()
		//get the global tracer
		tracer := opentracing.GlobalTracer()
		//this is where we start our span for this operation, this will be the parent for this method
		span := tracer.StartSpan("UUIDMiddleware")
		defer span.Finish()
		ctx = opentracing.ContextWithSpan(ctx, span)
		u2 := uuid.NewV4()
		/*if err != nil {
			respondError(w, http.StatusInternalServerError, err)
		}*/

		ctx = context.WithValue(ctx, "correlationid", u2.String())
		r = r.WithContext(ctx)
		next(w, r, ps)
	}
}
