package event

import (
	"errors"

	"gorm.io/gorm"
)

var (
	ErrNoRecord  error = errors.New("record not found")
	ErrNoRecords error = errors.New("records not found with filter")
)

type EventRepo struct {
	conn *gorm.DB
}

func NewRepo(conn *gorm.DB) *EventRepo {
	return &EventRepo{conn}
}

func (r *EventRepo) WithID(id int) (*Event, error) {
	event := Event{}
	err := r.conn.Where("id = ?", id).Find(&event).Error
	if err != nil {
		return &event, err
	}

	if event.ID == 0 {
		return &event, ErrNoRecord
	}

	return &event, err
}

func (r *EventRepo) Store(event *Event) error {
	return r.conn.Create(event).Error
}

func (r *EventRepo) Update(event *Event) error {
	return r.conn.Save(event).Error
}

func (r *EventRepo) Delete(event *Event) error {
	return r.conn.Delete(event).Error
}
