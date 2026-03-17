package storage

import (
	"context"
	"database/sql"
	"fmt"

	"forum/internal/models"
	"forum/internal/projfs"
)

func InitSQLite3DB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", string(projfs.StorageDir().Join("database.db")))
	if err != nil {
		return nil, err
	}

	query := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT NOT NULL UNIQUE,
		password_hash TEXT NOT NULL
	);
	CREATE TABLE IF NOT EXISTS sessions (
		token TEXT PRIMARY KEY,
		user_id INTEGER NOT NULL,
		FOREIGN KEY(user_id) REFERENCES users(id)
	);`

	_, err = db.Exec(query)
	return db, err
}

type UserRepo interface {
	// User || ErrUserNotExist
	ByLogin(ctx context.Context, login string) (*models.User, error)
	Save(ctx context.Context, user models.User) error
}

func NewSQLite3UserRepo(db *sql.DB) *sqlite3UserRepo {
	return &sqlite3UserRepo{db}
}

type sqlite3UserRepo struct {
	db *sql.DB
}

func (r sqlite3UserRepo) ByLogin(ctx context.Context, login string) (*models.User, error) {
	return nil, nil
}

func (r sqlite3UserRepo) Save(ctx context.Context, user models.User) error {
	return nil
}

var (
	ErrUserExists   = fmt.Errorf("user already exists")
	ErrUserNotExist = fmt.Errorf("user not exist")
)
