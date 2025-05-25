package service

import "crud-without-db/internal/domain"

//go:generate mockgen -source=user.go -destination=mocks/mock.go

type UsersRepo interface {
	Create(user domain.User) error
	GetByID(id int64) (domain.User, error)
	GetAll() ([]domain.User, error)
	Delete(id int64) error
	Update(id int64, inp domain.User) error
}

type Users struct {
	repo UsersRepo
}

func NewUsers(repo UsersRepo) *Users {
	return &Users{repo: repo}
}

func (s *Users) Create(user domain.User) error {
	return s.repo.Create(user)
}

func (s *Users) GetByID(id int64) (domain.User, error) {
	return s.repo.GetByID(id)
}

func (s *Users) GetAll() ([]domain.User, error) {
	return s.repo.GetAll()
}

func (s *Users) Delete(id int64) error {
	return s.repo.Delete(id)
}

func (s *Users) Update(id int64, inp domain.User) error {
	return s.repo.Update(id, inp)
}
