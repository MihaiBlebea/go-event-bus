package bus

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/MihaiBlebea/go-event-bus/project"
)

type TriggerRequest struct {
	Token   string `json:"token"`
	Event   string `json:"event"`
	Payload string `json:"payload"`
}

type TriggerResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}

func TriggerHandler(s Service, p project.Service) http.Handler {
	validate := func(r *http.Request) (*TriggerRequest, error) {
		request := TriggerRequest{}

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
		response := TriggerResponse{}

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

		if err := s.HandleIncomingEvent(
			projectID,
			request.Event,
			request.Payload); err != nil {

			response.Message = err.Error()
			sendResponse(w, response, http.StatusBadRequest)
			return
		}

		response.Success = true

		sendResponse(w, response, http.StatusOK)
	})
}
