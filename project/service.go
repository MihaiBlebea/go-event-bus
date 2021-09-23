package project

import "errors"

type Service interface {
	Create(name string) (string, error)
	Projects() ([]Project, error)
	ParseToken(token string) (int, error)
}

type service struct {
	projectRepo *ProjectRepo
}

func NewService(projectRepo *ProjectRepo) Service {
	return &service{projectRepo}
}

func (s *service) Create(name string) (string, error) {
	// check if a project with same slug exists
	slug := newName(name).toSlug()
	if _, err := s.projectRepo.WithSlug(slug); err == nil {
		return "", errors.New("project already exists")
	}

	project := New(name)
	if err := s.projectRepo.Store(project); err != nil {
		return "", err
	}

	return project.Token, nil
}

func (s *service) Projects() ([]Project, error) {
	return s.projectRepo.All()
}

func (s *service) ParseToken(token string) (int, error) {
	project, err := s.projectRepo.WithToken(token)
	if err != nil {
		return 0, err
	}

	return project.ID, nil
}
