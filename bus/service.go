package bus

import (
	"errors"

	"github.com/MihaiBlebea/go-event-bus/bus/event"
	"github.com/MihaiBlebea/go-event-bus/bus/sent"
	"github.com/MihaiBlebea/go-event-bus/bus/subscriber"
	"gorm.io/gorm"
)

type Service interface {
	Subscribe(projectID int, eventName, handlerUrl string) error
	Unsubscribe(projectID int, eventName string) error
	GetProjectSubscribers(projectID int) ([]subscriber.Subscriber, error)
	GetProcessedEvents(projectID int, pagination struct{ PerPage, Page int }) ([]sent.Sent, error)
	HandleIncomingEvent(projectID int, eventName, payload string) error
}

type service struct {
	eventRepo      *event.EventRepo
	subscriberRepo *subscriber.SubscriberRepo
	sentRepo       *sent.SentRepo
}

func NewService(conn *gorm.DB) Service {
	return &service{
		eventRepo:      event.NewRepo(conn),
		subscriberRepo: subscriber.NewRepo(conn),
		sentRepo:       sent.NewRepo(conn),
	}
}

func (s *service) Subscribe(projectID int, eventName, handlerUrl string) error {
	if _, err := s.subscriberRepo.WithEventName(eventName); err == nil {
		return errors.New("subscriber already exists")
	}

	sub := subscriber.New(projectID, eventName, handlerUrl, true)
	if err := s.subscriberRepo.Store(sub); err != nil {
		return err
	}

	return nil
}

func (s *service) Unsubscribe(projectID int, eventName string) error {
	sub, err := s.subscriberRepo.WithEventName(eventName)
	if err != nil {
		return err
	}

	if sub.ProjectID != projectID {
		return errors.New("resource not own")
	}

	return s.subscriberRepo.Delete(sub)
}

func (s *service) GetProjectSubscribers(projectID int) ([]subscriber.Subscriber, error) {
	return s.subscriberRepo.WithProjectID(projectID)
}

func (s *service) GetProcessedEvents(
	projectID int,
	paginate struct{ PerPage, Page int }) ([]sent.Sent, error) {

	offset := (paginate.Page - 1) * paginate.PerPage
	return s.sentRepo.WithProjectIDPaginated(projectID, offset, paginate.PerPage)
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
		if err := s.post(sub.HandlerUrl, payload); err != nil {
			s.sentEventFailed(projectID, &sub, eventName, err.Error())
			return err
		}

		s.sentEventSuccess(projectID, &sub, eventName)
	}

	return nil
}
