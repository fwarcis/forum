package models

import (
	"time"

	"forum/internal/common/system/coreuuid"
)

type Session struct {
	id        coreuuid.UUID
	expiresAt time.Time
}

type IDGenerator interface {
	Generate() (*coreuuid.UUID, error)
}

func newSession(idGener IDGenerator) (*Session, error) {
	id, err := idGener.Generate()
	if err != nil {
		return nil, err
	}

	return &Session{
		id:        *id,
		expiresAt: time.Now().Add(1 * time.Hour),
	}, nil
}

func (s Session) ID() coreuuid.UUID {
	return s.id
}

func (s Session) IsExpired() bool {
	return s.expiresAt.Before(time.Now())
}

func (s Session) ExpiresAt() time.Time {
	return s.expiresAt
}
