package apigateway

import (
	"encoding/json"
	"errors"
	"github.com/julienschmidt/httprouter"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"io"
	"net/http"
)

//using http router, register func will do the routing path registration
func (e Endpoints) Register(r *httprouter.Router) {
	//curl -XPOST localhost:8089/v1/login -d '{"email":"ff@ff.ff","password":"Aa111111"}'
	r.Handle("POST", "/v1/login", e.HandleLoginPost)
	//curl -XPOST localhost:8089/v1/register -d '{"email":"ff@ff.ffnew","password":"Aa111111", "firstname":"Farhad","lastname":"Farahi"}'
	r.Handle("POST", "/v1/register", e.HandleRegisterPost)
	r.Handler("GET", "/metrics", promhttp.Handler())
}

//each method needs a http handler handlers are registered in the register func
func (e Endpoints) HandleLoginPost(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	decodedLoginReq, err := decodeLoginRequest(e.Ctx, r)
	if err != nil {
		if err == io.EOF {
			respondError(w, http.StatusBadRequest, err)
			return
		}
		if err.Error() == "no rows in result set" {
			respondError(w, http.StatusBadRequest, errors.New("incorrect email or password"))
			return
		}
		respondError(w, 500, err)
		return
	}
	resp, err := e.LoginEndpoint(e.Ctx, decodedLoginReq.(loginRequest))
	if err != nil {
		respondError(w, 500, err)
		return
	}
	respondSuccess(w, resp.(loginResponse))
}

//each method needs a http handler handlers are registered in the register func
func (e Endpoints) HandleRegisterPost(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
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

// respondError in some canonical format.
func respondError(w http.ResponseWriter, code int, err error) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error":       err,
		"status_code": code,
		"status_text": http.StatusText(code),
	})
}

// respondSuccess in some canonical format.
func respondSuccess(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(data)
}
