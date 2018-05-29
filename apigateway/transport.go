package apigateway

import (
	"context"
	"encoding/json"
	"net/http"
)

//decode function for each method
func decodeLoginRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req loginRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, err
	}
	return req, nil
}

//single encode func for all
func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}
