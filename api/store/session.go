package store

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type Session struct {
	SessionId  uuid.UUID `json:"session_id"`
	UserId     int64     `json:"user_id"`
	IPAddress  string    `json:"ip"`
	UserAgent  string    `json:"user_agent"`
	CreatedAt  time.Time `json:"created_at"`
	LastSeenAt time.Time `json:"last_seen_at"`
}

type SessionManager struct {
	rClient *redis.Client
	prefix  string
	ttl     time.Duration
}

func NewSessionManager(redisClient *redis.Client, ttl time.Duration) *SessionManager {
	return &SessionManager{
		rClient: redisClient,
		prefix:  "session",
		ttl:     ttl,
	}
}

func NewRedisClient(addr, password string, db int) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	return client
}

func (s *SessionManager) Set(ctx context.Context, key string, data any, ttl time.Duration) error {
	return s.rClient.Set(ctx, key, data, ttl).Err()
}

func (s *SessionManager) Get(ctx context.Context, key string) (string, error) {
	data, err := s.rClient.Get(ctx, key).Result()
	return data, err
}

func (s *SessionManager) Del(ctx context.Context, key string) error {
	return s.rClient.Del(ctx, key).Err()
}

func (s *SessionManager) CreateSession(ctx context.Context, userId int64, ip, userAgent string) (*Session, error) {
	sessionId := uuid.New()

	sess := &Session{
		SessionId:  sessionId,
		UserId:     userId,
		IPAddress:  ip,
		UserAgent:  userAgent,
		CreatedAt:  time.Now(),
		LastSeenAt: time.Now(),
	}

	// Marshal data for storing in db
	data, err := json.Marshal(sess)
	if err != nil {
		return nil, err
	}

	key := fmt.Sprintf("%s_%s", s.prefix, sessionId.String())

	err = s.Set(ctx, key, data, s.ttl)
	if err != nil {
		return nil, err
	}

	return sess, nil
}

func (s *SessionManager) GetSession(ctx context.Context, sessionId uuid.UUID) (*Session, error) {
	key := fmt.Sprintf("%s_%s", s.prefix, sessionId.String())
	data, err := s.Get(ctx, key)
	if err != nil {
		return nil, err
	}

	var session Session
	if err := json.Unmarshal([]byte(data), &session); err != nil {
		return nil, err
	}

	return &session, nil
}

func (s *SessionManager) DeleteSession(ctx context.Context, sessionId uuid.UUID) error {
	key := fmt.Sprintf("%s_%s", s.prefix, sessionId.String())
	return s.Del(ctx, key)
}
