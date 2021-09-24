package project

import (
	"net/http"
	"time"
)

type ListResponse struct {
	Projects []ProjectResponse `json:"projects,omitempty"`
	Success  bool              `json:"success"`
	Message  string            `json:"message,omitempty"`
}

type ProjectResponse struct {
	Name    string    `json:"name"`
	Slug    string    `json:"slug"`
	Token   string    `json:"token"`
	Created time.Time `json:"created"`
}

func ListHandler(s Service) http.Handler {
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
			response.Projects = append(response.Projects, ProjectResponse{
				Name:    proj.Name,
				Slug:    proj.Slug,
				Token:   proj.Token,
				Created: proj.CreatedAt,
			})
		}

		sendResponse(w, response, http.StatusOK)
	})
}
