package store

import (
	"context"
	"database/sql"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Store struct {
	Users interface {
		RegisterUser(ctx context.Context, user *User) error
		GetByEmail(ctx context.Context, email string) (*User, error)
		GetById(ctx context.Context, userId int64) (*User, error)
	}
	Session interface {
		Set(ctx context.Context, key string, data any, ttl time.Duration) error
		Get(ctx context.Context, key string) (string, error)
		Del(ctx context.Context, key string) error
		CreateSession(ctx context.Context, userId int64, ip, userAgent string) (*Session, error)
		GetSession(ctx context.Context, sessionId uuid.UUID) (*Session, error)
		DeleteSession(ctx context.Context, sessionId uuid.UUID) error
	}
	Auth interface {
		GenerateToken(claims jwt.Claims) (string, error)
		ValidateToken(token string) (*jwt.Token, error)
	}
}

func NewStore(db *sql.DB, sessionMgr *SessionManager, auth *JWTAuthenticator) *Store {
	return &Store{
		Users:   &UserStore{db},
		Session: sessionMgr,
		Auth:    auth,
	}
}
