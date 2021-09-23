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
	Created    time.Time `json:"created"`
	Updated    time.Time `json:"updated"`
}

func New(projectID int, eventName, handlerUrl string, active bool) *Subscriber {
	return &Subscriber{
		ProjectID:  projectID,
		EventName:  eventName,
		HandlerUrl: handlerUrl,
		Active:     active,
	}
}
