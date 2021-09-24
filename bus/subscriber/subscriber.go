package subscriber

import (
	"time"
)

type Subscriber struct {
	ID         int       `json:"id"`
	ProjectID  int       `json:"project_id"`
	EventName  string    `json:"event_name"`
	HandlerUrl string    `json:"handler_url"`
	Active     bool      `json:"active"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func New(projectID int, eventName, handlerUrl string, active bool) *Subscriber {
	return &Subscriber{
		ProjectID:  projectID,
		EventName:  eventName,
		HandlerUrl: handlerUrl,
		Active:     active,
	}
}
