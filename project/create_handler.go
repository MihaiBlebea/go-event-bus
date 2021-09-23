package project

import (
	"encoding/json"
	"errors"
	"net/http"
)

type CreateRequest struct {
	Name string `json:"name"`
}

type CreateResponse struct {
	Token   string `json:"token,omitempty"`
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}

func CreateHandler(s Service) http.Handler {
	validate := func(r *http.Request) (*CreateRequest, error) {
		request := CreateRequest{}

		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			return &request, err
		}

		if request.Name == "" {
			return &request, errors.New("invalid request param name")
		}

		return &request, nil
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := CreateResponse{}

		request, err := validate(r)
		if err != nil {
			response.Message = err.Error()
			sendResponse(w, response, http.StatusBadRequest)
			return
		}

		token, err := s.Create(request.Name)
		if err != nil {
			response.Message = err.Error()
			sendResponse(w, response, http.StatusBadRequest)
			return
		}

		response.Success = true
		response.Token = token

		sendResponse(w, response, http.StatusOK)
	})
}

func sendResponse(w http.ResponseWriter, resp interface{}, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)

	b, _ := json.Marshal(resp)

	w.Write(b)
}
