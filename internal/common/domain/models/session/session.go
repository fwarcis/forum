package session

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type Session struct {
	id        uuid.UUID
	expiresAt time.Time
}

var ErrOnGenSessionID = errors.New("cannot generate session id")

func New() (*Session, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return nil, ErrOnGenSessionID
	}

	return &Session{
		id:        id,
		expiresAt: time.Now().Add(1 * time.Hour),
	}, nil
}

func (s Session) ID() uuid.UUID {
	return s.id
}

var ErrSessionExpired = errors.New("session expired")

func (s Session) Expired() error {
	if s.expiresAt.Before(time.Now()) {
		return ErrSessionExpired
	}
	return nil
}

func (s Session) ExpiresAt() time.Time {
	return s.expiresAt
}
