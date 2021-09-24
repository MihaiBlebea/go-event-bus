package bus

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/MihaiBlebea/go-event-bus/bus/event"
	"github.com/MihaiBlebea/go-event-bus/bus/subscriber"
	"gorm.io/gorm"
)

type Service interface {
	AddSubscriber(projectID int, eventName, handlerUrl string) error
	GetProjectSubscribers(projectID int) ([]subscriber.Subscriber, error)
	HandleIncomingEvent(projectID int, eventName, payload string) error
}

type service struct {
	eventRepo      *event.EventRepo
	subscriberRepo *subscriber.SubscriberRepo
}

func NewService(conn *gorm.DB) Service {
	return &service{
		eventRepo:      event.NewRepo(conn),
		subscriberRepo: subscriber.NewRepo(conn),
	}
}

func (s *service) AddSubscriber(projectID int, eventName, handlerUrl string) error {
	if _, err := s.subscriberRepo.WithEventName(eventName); err == nil {
		return errors.New("subscriber already exists")
	}

	sub := subscriber.New(projectID, eventName, handlerUrl, true)
	if err := s.subscriberRepo.Store(sub); err != nil {
		return err
	}

	return nil
}

func (s *service) GetProjectSubscribers(projectID int) ([]subscriber.Subscriber, error) {
	return s.subscriberRepo.WithProjectID(projectID)
}

func (s *service) HandleIncomingEvent(projectID int, eventName, payload string) error {
	event := event.New(projectID, eventName, payload)
	if err := s.eventRepo.Store(event); err != nil {
		return err
	}

	subs, err := s.subscriberRepo.WithProjectID(projectID)
	if err != nil {
		return err
	}

	for _, sub := range subs {
		if err := post(sub.HandlerUrl, payload); err != nil {
			return err
		}
	}

	return nil
}

func post(endpoint string, payload interface{}) error {
	b, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(
		"POST",
		endpoint,
		bytes.NewBuffer(b),
	)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return errors.New("request status is not 200")
	}

	return nil
}
