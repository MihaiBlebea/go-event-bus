package sent

import (
	"errors"

	"gorm.io/gorm"
)

var (
	ErrNoRecord  error = errors.New("record not found")
	ErrNoRecords error = errors.New("records not found with filter")
)

type SentRepo struct {
	conn *gorm.DB
}

func NewRepo(conn *gorm.DB) *SentRepo {
	return &SentRepo{conn}
}

func (r *SentRepo) WithID(id int) (*Sent, error) {
	sent := Sent{}
	err := r.conn.Where("id = ?", id).Find(&sent).Error
	if err != nil {
		return &sent, err
	}

	if sent.ID == 0 {
		return &sent, ErrNoRecord
	}

	return &sent, err
}

func (r *SentRepo) WithProjectID(id int) ([]Sent, error) {
	sents := []Sent{}
	err := r.conn.Where("project_id = ?", id).Find(&sents).Error
	if err != nil {
		return sents, err
	}

	return sents, err
}

func (r *SentRepo) WithProjectIDPaginated(id, offset, pageSize int) ([]Sent, error) {
	sents := []Sent{}
	err := r.conn.Where("project_id = ?", id).Offset(offset).Limit(pageSize).Find(&sents).Error
	if err != nil {
		return sents, err
	}

	return sents, err
}

func (r *SentRepo) Store(sent *Sent) error {
	return r.conn.Create(sent).Error
}

func (r *SentRepo) Update(sent *Sent) error {
	return r.conn.Save(sent).Error
}

func (r *SentRepo) Delete(sent *Sent) error {
	return r.conn.Delete(sent).Error
}
