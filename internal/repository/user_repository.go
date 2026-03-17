package repository

import (
	"database/sql"
	"time"

	"template.go/internal/model"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// UserWithPassword — helper struct to handle authentication logic
type UserWithPassword struct {
	model.User
	Password string // Matches 'password_hash' in DB
}

func (r *UserRepository) Create(name, email, passwordHash string) (*model.User, error) {
	query := `
		INSERT INTO users (name, email, password, created_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`
	user := &model.User{
		Name:      name,
		Email:     email,
		CreatedAt: time.Now(),
	}

	err := r.db.QueryRow(query, name, email, passwordHash, user.CreatedAt).Scan(&user.ID)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) FindByEmail(email string) (*UserWithPassword, error) {
	query := `SELECT id, name, email, password, created_at FROM users WHERE email = $1`
	u := &UserWithPassword{}

	err := r.db.QueryRow(query, email).Scan(
		&u.ID, &u.Name, &u.Email, &u.Password, &u.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (r *UserRepository) FindByID(id int64) (*model.User, error) {
	query := `SELECT id, name, email, created_at FROM users WHERE id = $1`
	u := &model.User{}

	err := r.db.QueryRow(query, id).Scan(&u.ID, &u.Name, &u.Email, &u.CreatedAt)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (r *UserRepository) EmailExists(email string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)`
	err := r.db.QueryRow(query, email).Scan(&exists)
	return exists, err
}
