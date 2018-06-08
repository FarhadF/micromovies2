package apigateway

import (
	"context"
	"github.com/julienschmidt/httprouter"
	"github.com/satori/go.uuid"
	"net/http"
)

/*type uUIDMiddleware struct {
	ctx  context.Context
	next *httprouter.Router
}

func NewUUIDMiddleware(ctx context.Context, next *httprouter.Router) *uUIDMiddleware {
	return &uUIDMiddleware{ctx, next}
}*/

/*func (e *uUIDMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println("here")
	// We can modify the request here
	u2, err := uuid.NewV4()
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
	}
	e.Ctx = context.WithValue(e.Ctx, "correlationid",u2)
	r.WithContext(e.Ctx)
	e.Next.ServeHTTP(w, r)
	// We can modify the response here
}*/

func UUIDMiddleware(next httprouter.Handle) httprouter.Handle {

	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		u2, err := uuid.NewV4()
		if err != nil {
			respondError(w, http.StatusInternalServerError, err)
		}
		ctx := r.Context()
		ctx = context.WithValue(ctx, "correlationid", u2)
		//fmt.Println(ctx)
		r = r.WithContext(ctx)
		next(w, r, ps)
	}
}
