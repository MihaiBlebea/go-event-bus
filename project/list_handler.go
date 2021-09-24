package project

import (
	"net/http"
	"time"

	"github.com/MihaiBlebea/go-event-bus/bus/subscriber"
)

type BusService interface {
	GetProjectSubscribers(projectID int) ([]subscriber.Subscriber, error)
}

type ListResponse struct {
	Projects []ProjectResponse `json:"projects,omitempty"`
	Success  bool              `json:"success"`
	Message  string            `json:"message,omitempty"`
}

type ProjectResponse struct {
	Name        string               `json:"name"`
	Slug        string               `json:"slug"`
	Token       string               `json:"token"`
	Subscribers []SubscriberResponse `json:"subscribers"`
	Created     time.Time            `json:"created"`
}

type SubscriberResponse struct {
	Event string `json:"event"`
	URL   string `json:"url"`
}

func ListHandler(s Service, b BusService) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := ListResponse{}

		projects, err := s.Projects()
		if err != nil {
			response.Message = err.Error()
			sendResponse(w, response, http.StatusBadRequest)
			return
		}

		response.Success = true
		for _, proj := range projects {
			subs, err := b.GetProjectSubscribers(proj.ID)
			if err != nil {
				response.Message = err.Error()
				sendResponse(w, response, http.StatusBadRequest)
				return
			}

			var subsResp []SubscriberResponse
			for _, sub := range subs {
				subsResp = append(subsResp, SubscriberResponse{
					Event: sub.EventName,
					URL:   sub.HandlerUrl,
				})
			}

			response.Projects = append(response.Projects, ProjectResponse{
				Name:        proj.Name,
				Slug:        proj.Slug,
				Token:       proj.Token,
				Subscribers: subsResp,
				Created:     proj.CreatedAt,
			})
		}

		sendResponse(w, response, http.StatusOK)
	})
}
