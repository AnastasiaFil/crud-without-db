package service

import (
	"crud-without-db/internal/domain"
)

type UsersRepository interface {
	Create(user domain.User) error
	GetByID(id int64) (domain.User, error)
	GetAll() ([]domain.User, error)
	Delete(id int64) error
	Update(id int64, inp domain.User) error
}

type Users struct {
	repo UsersRepository
}

func NewUsers(repo UsersRepository) *Users {
	return &Users{
		repo: repo,
	}
}

func (b *Users) Create(user domain.User) error {
	return b.repo.Create(user)
}

func (b *Users) GetByID(id int64) (domain.User, error) {
	return b.repo.GetByID(id)
}

func (b *Users) GetAll() ([]domain.User, error) {
	return b.repo.GetAll()
}

func (b *Users) Delete(id int64) error {
	return b.repo.Delete(id)
}

func (b *Users) Update(id int64, inp domain.User) error {
	return b.repo.Update(id, inp)
}
