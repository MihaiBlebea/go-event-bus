package event

import (
	"time"
)

type Event struct {
	ID        int       `json:"id"`
	ProjectID int       `json:"project_id"`
	EventName string    `json:"event_name"`
	Payload   string    `json:"payload"`
	Created   time.Time `json:"created"`
	Updated   time.Time `json:"updated"`
}

func New(projectID int, eventName string, payload string) *Event {
	return &Event{
		ProjectID: projectID,
		EventName: eventName,
		Payload:   payload,
	}
}
