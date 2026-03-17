package servs

import (
	"context"
	"errors"
	"time"

	"forum/internal/models"
	"forum/internal/storage"
)

type UserService struct {
	repo storage.UserRepo
}

func NewUserService(repo storage.UserRepo) *UserService {
	return &UserService{repo}
}

func (s UserService) Register(
	ctx context.Context,
	login string,
	email string,
	password []byte,
) (*SessionData, error) {
	_, err := s.repo.ByLogin(ctx, login)
	if !errors.Is(err, storage.ErrUserNotExist) {
		return nil, err
	}

	user, err := models.NewUser(login, email, password)
	if err != nil {
		return nil, err
	}

	err = s.repo.Save(ctx, user)
	if err != nil {
		return nil, err
	}
	sessnData := newSessionData(user.Session())
	return &sessnData, nil
}

func (s UserService) LogIn(
	ctx context.Context,
	login string,
	password []byte,
) (*SessionData, error) {
	user, err := s.repo.ByLogin(ctx, login)
	if err != nil {
		return nil, err
	}
	err = user.CheckPassword(password)
	if err != nil {
		return nil, err
	}

	if !user.Session().IsExpired() {
		sessnData := newSessionData(user.Session())
		return &sessnData, nil
	}
	err = user.GenerateNewSession()
	if err != nil {
		return nil, err
	}
	sessnData := newSessionData(user.Session())
	return &sessnData, nil
}

func newSessionData(sessn models.Session) SessionData {
	return SessionData{
		ID:        sessn.ID().String(),
		ExpiresAt: sessn.ExpiresAt(),
	}
}

type SessionData struct {
	ID        string
	ExpiresAt time.Time
}
