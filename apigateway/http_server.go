package apigateway

import (
	"encoding/json"
	"errors"
	"github.com/julienschmidt/httprouter"
	"github.com/opentracing/opentracing-go"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"io"
	"net/http"
	"strings"
)

//todo: api documentation
//using http router, register func will do the routing path registration
func (e Endpoints) Register(r *httprouter.Router) {
	//curl -XPOST localhost:8089/v1/login -d '{"email":"ff@ff.ff","password":"Aa111111"}'
	r.Handle("POST", "/v1/login", UUIDMiddleware(e.HandleLoginPost))
	// swagger:route POST /login login users login
	// Authenticates user
	// responses:
	//  200: loginResponse
	//  400: loginResponse
	//curl -XPOST localhost:8089/v1/register -d '{"email":"ff@ff.ffnew","password":"Aa111111", "firstname":"Farhad","lastname":"Farahi"}'
	r.Handle("POST", "/v1/register", e.HandleRegisterPost)
	//curl -XPOST localhost:8089/v1/changepassword -d '{"email":"ff@ff.ff","currentpassword":"Aa111111","newpassword":"Aa123"}' --header "Authorization: Bearer ..."
	r.Handle("POST", "/v1/changepassword", UUIDMiddleware(e.HandleChangePasswordPost))
	r.Handler("GET", "/metrics", promhttp.Handler())
}

//each method needs a http handler handlers are registered in the register func
func (e Endpoints) HandleLoginPost(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	//take out http request context that we put in at auth middleware and put it in go-kit endpoint context
	e.Ctx = r.Context()
	if span := opentracing.SpanFromContext(e.Ctx); span != nil {
		span := span.Tracer().StartSpan("HandleLoginPost", opentracing.ChildOf(span.Context()))
		defer span.Finish()
		e.Ctx = opentracing.ContextWithSpan(e.Ctx, span)
	}
	decodedLoginReq, err := decodeLoginRequest(e.Ctx, r)
	if err != nil {
		respondError(w, http.StatusBadRequest, errors.New("incorrect email or password"))
		return
	}
	resp, err := e.LoginEndpoint(e.Ctx, decodedLoginReq.(loginRequest))
	res := resp.(loginResponse)

	if err != nil {
		respondError(w, 500, err)
		return
	}
	if res.Err != "" {
		if strings.EqualFold(res.Err, "email or password incorrect") {
			respondError(w, http.StatusBadRequest, errors.New("incorrect email or password"))
			return
		}
		respondError(w, 500, errors.New(res.Err))
		return
	}
	respondSuccess(w, res)
}

//each method needs a http handler handlers are registered in the register func
func (e Endpoints) HandleRegisterPost(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	//take out http request context that we put in at auth middleware and put it in go-kit endpoint context
	e.Ctx = r.Context()
	decodedRegisterReq, err := decodeRegisterRequest(e.Ctx, r)
	if err != nil {
		if err == io.EOF {
			respondError(w, http.StatusBadRequest, err)
			return
		}
		respondError(w, 500, err)
		return
	}
	resp, err := e.RegisterEndpoint(e.Ctx, decodedRegisterReq.(registerRequest))
	if err != nil {
		respondError(w, 500, err)
		return
	}
	respondSuccess(w, resp.(registerResponse))
}

//each method needs a http handler handlers are registered in the register func
func (e Endpoints) HandleChangePasswordPost(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	//take out http request context that we put in at auth middleware and put it in go-kit endpoint context
	e.Ctx = r.Context()
	decodedChangePasswordReq, err := decodeChangePasswordRequest(e.Ctx, r)
	if err != nil {
		if err == io.EOF {
			respondError(w, http.StatusBadRequest, err)
			return
		}
		respondError(w, 500, err)
		return
	}
	resp, err := e.ChangePasswordEndpoint(e.Ctx, decodedChangePasswordReq.(changePasswordRequest))
	if err != nil {
		respondError(w, 500, err)
		return
	}
	respondSuccess(w, resp.(changePasswordResponse))
}

// respondError in some canonical format.
func respondError(w http.ResponseWriter, code int, err error) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error":       err.Error(),
		"status_code": code,
		"status_text": http.StatusText(code),
	})
}

// respondSuccess in some canonical format.
//todo: check encoding errors
func respondSuccess(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(data)
}
