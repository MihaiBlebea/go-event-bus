package bus

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/MihaiBlebea/go-event-bus/project"
)

type UnsubscribeRequest struct {
	Token string `json:"token"`
	Event string `json:"event"`
}

type UnsubscribeResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}

func UnsubscribeHandler(s Service, p project.Service) http.Handler {
	validate := func(r *http.Request) (*UnsubscribeRequest, error) {
		request := UnsubscribeRequest{}

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

		return &request, nil
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := UnsubscribeResponse{}

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

		if err := s.Unsubscribe(projectID, request.Event); err != nil {
			response.Message = err.Error()
			sendResponse(w, response, http.StatusBadRequest)
			return
		}

		response.Success = true

		sendResponse(w, response, http.StatusOK)
	})
}
