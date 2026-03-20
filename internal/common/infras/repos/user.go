package repos

import (
	"context"
	"database/sql"

	"forum/internal/common/models"
	"forum/internal/common/system/projfs"
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

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{db}
}

type UserRepo struct {
	db *sql.DB
}

func (r UserRepo) ByLogin(ctx context.Context, login string) (*models.User, error) {
	return nil, nil
}

func (r UserRepo) Save(ctx context.Context, user models.User) error {
	return nil
}

func (r UserRepo) Commit(ctx context.Context) error {
	return nil
}
