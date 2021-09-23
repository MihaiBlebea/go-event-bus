package subscriber

import (
	"errors"

	"gorm.io/gorm"
)

var (
	ErrNoRecord  error = errors.New("record not found")
	ErrNoRecords error = errors.New("records not found with filter")
)

type SubscriberRepo struct {
	conn *gorm.DB
}

func NewRepo(conn *gorm.DB) *SubscriberRepo {
	return &SubscriberRepo{conn}
}

func (r *SubscriberRepo) WithID(id int) (*Subscriber, error) {
	subscriber := Subscriber{}
	err := r.conn.Where("id = ?", id).Find(&subscriber).Error
	if err != nil {
		return &subscriber, err
	}

	if subscriber.ID == 0 {
		return &subscriber, ErrNoRecord
	}

	return &subscriber, err
}

func (r *SubscriberRepo) WithProjectID(id int) ([]Subscriber, error) {
	subscribers := []Subscriber{}
	err := r.conn.Where("project_id = ?", id).Find(&subscribers).Error
	if err != nil {
		return subscribers, err
	}

	return subscribers, err
}

func (r *SubscriberRepo) Store(subscriber *Subscriber) error {
	return r.conn.Create(subscriber).Error
}

func (r *SubscriberRepo) Update(subscriber *Subscriber) error {
	return r.conn.Save(subscriber).Error
}

func (r *SubscriberRepo) Delete(subscriber *Subscriber) error {
	return r.conn.Delete(subscriber).Error
}
