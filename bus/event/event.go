package event

import (
	"time"
)

type Event struct {
	ID        int       `json:"id"`
	ProjectID int       `json:"project_id"`
	EventName string    `json:"event_name"`
	Payload   string    `json:"payload"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func New(projectID int, eventName string, payload string) *Event {
	return &Event{
		ProjectID: projectID,
		EventName: eventName,
		Payload:   payload,
	}
}
