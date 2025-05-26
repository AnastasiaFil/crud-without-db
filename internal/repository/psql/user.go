package psql

import (
	"crud-without-db/internal/domain"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

type Users struct {
	db *sql.DB
}

func NewUsers(db *sql.DB) *Users {
	return &Users{db: db}
}

func (r *Users) Create(user domain.User) error {
	query := `
		INSERT INTO users (name, age, sex) 
		VALUES ($1, $2, $3) 
		RETURNING id`

	err := r.db.QueryRow(query, user.Name, user.Age, user.Sex).Scan(&user.ID)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

func (r *Users) GetByID(id int64) (domain.User, error) {
	var user domain.User
	query := `SELECT id, name, age, sex FROM users WHERE id = $1`

	err := r.db.QueryRow(query, id).Scan(&user.ID, &user.Name, &user.Age, &user.Sex)
	if err != nil {
		if err == sql.ErrNoRows {
			return user, domain.ErrUserNotFound
		}
		return user, fmt.Errorf("failed to get user by id: %w", err)
	}

	return user, nil
}

func (r *Users) GetAll() ([]domain.User, error) {
	query := `SELECT id, name, age, sex FROM users ORDER BY id`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get all users: %w", err)
	}
	defer rows.Close()

	var users []domain.User
	for rows.Next() {
		var user domain.User
		err := rows.Scan(&user.ID, &user.Name, &user.Age, &user.Sex)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return users, nil
}

func (r *Users) Update(id int64, user domain.User) error {
	query := `
		UPDATE users 
		SET name = $1, age = $2, sex = $3 
		WHERE id = $4`

	result, err := r.db.Exec(query, user.Name, user.Age, user.Sex, id)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return domain.ErrUserNotFound
	}

	return nil
}

func (r *Users) Delete(id int64) error {
	query := `DELETE FROM users WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return domain.ErrUserNotFound
	}

	return nil
}

// InitSchema creates the users table if it doesn't exist
func (r *Users) InitSchema() error {
	query := `
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			age INTEGER NOT NULL,
			sex VARCHAR(10) NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`

	_, err := r.db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create users table: %w", err)
	}

	return nil
}
