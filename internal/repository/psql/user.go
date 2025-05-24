package psql

import (
	"crud-without-db/internal/domain"
)

type Users struct {
	users  []domain.User
	nextID int64
}

func NewUsers() *Users {
	return &Users{
		users:  []domain.User{},
		nextID: 1,
	}
}

func (r *Users) Create(user domain.User) error {
	newUser := domain.User{
		ID:   r.nextID,
		Name: user.Name,
		Age:  user.Age,
		Sex:  user.Sex,
	}
	r.users = append(r.users, newUser)
	r.nextID++
	return nil
}

func (r *Users) GetByID(id int64) (domain.User, error) {
	for _, u := range r.users {
		if u.ID == id {
			return u, nil
		}
	}
	return domain.User{}, domain.ErrUserNotFound
}

func (r *Users) GetAll() ([]domain.User, error) {
	return r.users, nil
}

func (r *Users) Delete(id int64) error {
	for i, u := range r.users {
		if u.ID == id {
			r.users = append(r.users[:i], r.users[i+1:]...)
			return nil
		}
	}
	return nil // No error even if user is not found
}

func (r *Users) Update(id int64, user domain.User) error {
	for i, u := range r.users {
		if u.ID == id {
			updatedUser := u
			if user.Name != "" {
				updatedUser.Name = user.Name
			}
			if user.Age != 0 {
				updatedUser.Age = user.Age
			}
			if user.Sex != "" {
				updatedUser.Sex = user.Sex
			}
			r.users[i] = updatedUser
			return nil
		}
	}
	return domain.ErrUserNotFound
}
