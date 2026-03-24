package servs

import (
	"context"
	"errors"
	"time"

	"forum/internal/common/domain/models/session"
	"forum/internal/common/domain/models/user"

	"github.com/google/uuid"
)

type userRepo interface {
	ByLogin(ctx context.Context, login string) (*user.User, error)
	BySessionID(ctx context.Context, id uuid.UUID) (*user.User, error)
	Save(ctx context.Context, usr user.User) error
	Commit(ctx context.Context) error
}

type IAMService struct {
	repo       userRepo
	passHasher user.PasswordHasher
}

func NewIAMService(
	repo userRepo,
	passHasher user.PasswordHasher,
) *IAMService {
	return &IAMService{repo, passHasher}
}

func (s IAMService) Auth(
	ctx context.Context,
	sessionID uuid.UUID,
) (context.Context, error) {
	usr, err := s.repo.BySessionID(ctx, sessionID)
	if err != nil {
		return nil, err
	}
	err = usr.Session().Expired()
	if err != nil {
		return nil, err
	}
	return user.NewContext(ctx, *usr), nil
}

func (s IAMService) Register(
	ctx context.Context,
	login string,
	email string,
	password string,
) (*SessionData, error) {
	_, err := s.repo.ByLogin(ctx, login)
	if !errors.Is(err, user.ErrNotFound) {
		return nil, err
	}

	user, err := user.New(login, email, password, s.passHasher)
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

	sessionData := newSessionData(user.Session())
	return &sessionData, nil
}

func (s IAMService) LogIn(
	ctx context.Context,
	login string,
	password string,
) (*SessionData, error) {
	usr, err := s.repo.ByLogin(ctx, login)
	if err != nil {
		return nil, err
	}
	err = usr.VerifyPassword(password)
	if err != nil {
		return nil, err
	}

	if usr.Session().Expired() == nil {
		sessionData := newSessionData(usr.Session())
		return &sessionData, nil
	}
	err = usr.GenerateNewSession()
	if err != nil {
		return nil, err
	}

	err = s.repo.Save(ctx, *usr)
	if err != nil {
		return nil, err
	}
	err = s.repo.Commit(ctx)
	if err != nil {
		return nil, err
	}

	sessionData := newSessionData(usr.Session())
	return &sessionData, nil
}

func newSessionData(sess session.Session) SessionData {
	id := sess.ID()
	return SessionData{
		ID:        string(id[:]),
		ExpiresAt: sess.ExpiresAt(),
	}
}

type SessionData struct {
	ID        string
	ExpiresAt time.Time
}
