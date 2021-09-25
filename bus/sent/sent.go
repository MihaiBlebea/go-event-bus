package sent

import (
	"time"
)

type Sent struct {
	ID           int       `json:"id"`
	ProjectID    int       `json:"project_id"`
	SubscriberID int       `json:"subscriber_id"`
	Name         string    `json:"name"`
	Url          string    `json:"url"`
	Success      bool      `json:"success"`
	ErrorMessage string    `json:"error_message"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func New(
	projectID, subscriberID int,
	name, url string,
	errorMessage string) *Sent {

	var success bool
	if errorMessage == "" {
		success = true
	}

	return &Sent{
		ProjectID:    projectID,
		SubscriberID: subscriberID,
		Name:         name,
		Url:          url,
		Success:      success,
		ErrorMessage: errorMessage,
	}
}
