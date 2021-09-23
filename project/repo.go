package project

import (
	"errors"

	"gorm.io/gorm"
)

var (
	ErrNoRecord  error = errors.New("record not found")
	ErrNoRecords error = errors.New("records not found with filter")
)

type ProjectRepo struct {
	conn *gorm.DB
}

func NewRepo(conn *gorm.DB) *ProjectRepo {
	return &ProjectRepo{conn}
}

func (r *ProjectRepo) WithID(id int) (*Project, error) {
	project := Project{}
	err := r.conn.Where("id = ?", id).Find(&project).Error
	if err != nil {
		return &project, err
	}

	if project.ID == 0 {
		return &project, ErrNoRecord
	}

	return &project, err
}

func (r *ProjectRepo) WithToken(token string) (*Project, error) {
	project := Project{}
	err := r.conn.Where("token = ?", token).Find(&project).Error
	if err != nil {
		return &project, err
	}

	if project.ID == 0 {
		return &project, ErrNoRecord
	}

	return &project, err
}

func (r *ProjectRepo) WithSlug(slug string) (*Project, error) {
	project := Project{}
	err := r.conn.Where("slug = ?", slug).Find(&project).Error
	if err != nil {
		return &project, err
	}

	if project.ID == 0 {
		return &project, ErrNoRecord
	}

	return &project, err
}

func (r *ProjectRepo) Store(project *Project) error {
	return r.conn.Create(project).Error
}

func (r *ProjectRepo) Update(project *Project) error {
	return r.conn.Save(project).Error
}

func (r *ProjectRepo) Delete(project *Project) error {
	return r.conn.Delete(project).Error
}
