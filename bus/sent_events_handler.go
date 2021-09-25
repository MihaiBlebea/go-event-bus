package bus

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/MihaiBlebea/go-event-bus/project"
)

type SentEventsRequest struct {
	Token   string `json:"token"`
	PerPage int    `json:"per_page"`
	Page    int    `json:"page"`
}

type SentEventsResponse struct {
	Success bool            `json:"success"`
	Message string          `json:"message,omitempty"`
	Events  []EventResponse `json:"events,omitempty"`
}

type EventResponse struct {
	Name         string    `json:"name"`
	Url          string    `json:"url"`
	ErrorMessage string    `json:"error_message,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
}

func SentEventsHandler(s Service, p project.Service) http.Handler {
	validate := func(r *http.Request) (*SentEventsRequest, error) {
		params := r.URL.Query()

		token := params.Get("token")
		perPage := params.Get("per_page")
		page := params.Get("page")

		perPageInt, err := strconv.Atoi(perPage)
		if err != nil {
			return &SentEventsRequest{}, errors.New("invalid query param per_page")
		}

		pageInt, err := strconv.Atoi(page)
		if err != nil {
			return &SentEventsRequest{}, errors.New("invalid query param page")
		}

		request := SentEventsRequest{token, perPageInt, pageInt}

		if request.Token == "" {
			return &request, errors.New("invalid query param token")
		}

		if request.PerPage == 0 {
			return &request, errors.New("invalid query param per_page")
		}

		if request.Page == 0 {
			return &request, errors.New("invalid query param page")
		}

		return &request, nil
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := SentEventsResponse{}

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

		paginate := struct{ PerPage, Page int }{request.PerPage, request.Page}
		events, err := s.GetProcessedEvents(projectID, paginate)
		if err != nil {
			response.Message = err.Error()
			sendResponse(w, response, http.StatusBadRequest)
			return
		}

		for _, ev := range events {
			response.Events = append(response.Events, EventResponse{
				Name:      ev.Name,
				Url:       ev.Url,
				CreatedAt: ev.CreatedAt,
			})
		}
		response.Success = true

		sendResponse(w, response, http.StatusOK)
	})
}
