package project

import (
	"net/http"
)

type ListResponse struct {
	Projects []Project `json:"projects,omitempty"`
	Success  bool      `json:"success"`
	Message  string    `json:"message,omitempty"`
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
		response.Projects = projects

		sendResponse(w, response, http.StatusOK)
	})
}
