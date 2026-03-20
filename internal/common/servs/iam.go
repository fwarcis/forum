package servs

import (
	"context"
	"errors"
	"time"

	"forum/internal/common/models"
)

type UserRepo interface {
	ByLogin(ctx context.Context, login string) (*models.User, error)
	Save(ctx context.Context, user models.User) error
	Commit(ctx context.Context) error
}

type IAMService struct {
	repo       UserRepo
	passHasher models.PasswordHasher
	idGener    models.IDGenerator
}

func NewIAMService(
	repo UserRepo,
	passHasher models.PasswordHasher,
	idGener models.IDGenerator,
) *IAMService {
	return &IAMService{repo, passHasher, idGener}
}

func (s IAMService) CheckSession()

func (s IAMService) Register(
	ctx context.Context,
	login string,
	email string,
	password []byte,
) (*SessionData, error) {
	_, err := s.repo.ByLogin(ctx, login)
	if !errors.Is(err, models.ErrUserNotFound) {
		return nil, err
	}

	user, err := models.NewUser(login, email, password, s.passHasher, s.idGener)
	if err != nil {
		return nil, err
	}

	err = s.repo.Save(ctx, *user)
	if err != nil {
		return nil, err
	}
	err = s.repo.Commit(ctx)
	if err != nil {
		return nil, err
	}

	sessnData := newSessionData(user.Session())
	return &sessnData, nil
}

func (s IAMService) LogIn(
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

	err = s.repo.Save(ctx, *user)
	if err != nil {
		return nil, err
	}
	err = s.repo.Commit(ctx)
	if err != nil {
		return nil, err
	}

	sessnData := newSessionData(user.Session())
	return &sessnData, nil
}

func newSessionData(session models.Session) SessionData {
	id := session.ID()
	idstr := string(id[:])
	return SessionData{
		ID:        idstr,
		ExpiresAt: session.ExpiresAt(),
	}
}

type SessionData struct {
	ID        string
	ExpiresAt time.Time
}
