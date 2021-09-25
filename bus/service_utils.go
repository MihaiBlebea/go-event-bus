package bus

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/MihaiBlebea/go-event-bus/bus/sent"
	"github.com/MihaiBlebea/go-event-bus/bus/subscriber"
)

func (s *service) post(endpoint string, payload interface{}) error {
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

func (s *service) sentEventFailed(
	projectID int,
	sub *subscriber.Subscriber,
	event, message string) error {

	return s.sentRepo.Store(
		sent.New(projectID, sub.ID, event, sub.HandlerUrl, message),
	)
}

func (s *service) sentEventSuccess(
	projectID int,
	sub *subscriber.Subscriber,
	event string) error {

	return s.sentRepo.Store(
		sent.New(projectID, sub.ID, event, sub.HandlerUrl, ""),
	)
}
