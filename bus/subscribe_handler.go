package bus

import (
	"encoding/json"
	"errors"
	"net/http"
)

type ProjectService interface {
	ParseToken(token string) (int, error)
}

type SubscribeRequest struct {
	Token string `json:"token"`
	Event string `json:"event"`
	Url   string `json:"url"`
}

type SubscribeResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}

func SubscribeHandler(s Service, p ProjectService) http.Handler {
	validate := func(r *http.Request) (*SubscribeRequest, error) {
		request := SubscribeRequest{}

		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			return &request, err
		}

		if request.Token == "" {
			return &request, errors.New("invalid request param token")
		}

		if request.Event == "" {
			return &request, errors.New("invalid request param event")
		}

		if request.Url == "" {
			return &request, errors.New("invalid request param url")
		}

		return &request, nil
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := SubscribeResponse{}

		request, err := validate(r)
		if err != nil {
			response.Message = err.Error()
			sendResponse(w, response, http.StatusBadRequest)
			return
		}

		projectID, err := p.ParseToken(request.Token)
		if err != nil {
			response.Message = err.Error()
			sendResponse(w, response, http.StatusBadRequest)
			return
		}

		if err := s.AddSubscriber(projectID, request.Event, request.Url); err != nil {
			response.Message = err.Error()
			sendResponse(w, response, http.StatusBadRequest)
			return
		}

		response.Success = true

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
