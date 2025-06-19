package store

import (
	"context"
	"database/sql"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id        int64    `json:"id"`
	Email     string   `json:"email"`
	Password  password `json:"password"`
	CreatedAt string   `json:"created_at"`
}

type password struct {
	text *string
	hash []byte
}

func (p *password) Set(textPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(textPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	p.text = &textPassword
	p.hash = hash

	return nil
}

func (p *password) ComparePassword(textPassword string) error {
	return bcrypt.CompareHashAndPassword(p.hash, []byte(textPassword))
}

type UserStore struct {
	db *sql.DB
}

func (s *UserStore) RegisterUser(ctx context.Context, user *User) error {
	query := `
		INSERT INTO users (email, password) VALUES ($1, $2) RETURNING id, created_at
	`

	ctx, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()

	err := s.db.QueryRowContext(ctx, query, user.Email, user.Password.hash).Scan(
		&user.Id,
		&user.CreatedAt,
	)

	if err != nil {
		return err
	}

	return nil
}

func (s *UserStore) GetByEmail(ctx context.Context, email string) (*User, error) {
	query := `
		SELECT id, email, password, created_at
		FROM users
		WHERE email=$1
	`

	var user User

	err := s.db.QueryRowContext(ctx, query, email).Scan(
		&user.Id,
		&user.Email,
		&user.Password.hash,
		&user.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *UserStore) GetById(ctx context.Context, userId int64) (*User, error) {
	query := `
		SELECT id, email, created_at
		FROM users
		WHERE id=$1
	`

	var user User

	err := s.db.QueryRowContext(ctx, query, userId).Scan(
		&user.Id,
		&user.Email,
		&user.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}
