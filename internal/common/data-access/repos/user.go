package repos

import (
	"context"
	"database/sql"
	"os/user"
)

func NewUserRepo(db *sql.DB) *userRepo {
	return &userRepo{db}
}

type userRepo struct {
	db *sql.DB
}

func (r userRepo) ByLogin(ctx context.Context, login string) (*user.User, error) {
	return nil, nil
}

func (r userRepo) Save(ctx context.Context, usr user.User) error {
	return nil
}

func (r userRepo) Commit(ctx context.Context) error {
	return nil
}
