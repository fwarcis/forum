package main

import (
	"database/sql"
	"log/slog"
	"net/http"
	"os"

	"forum/internal/common/infras/passhash"
	"forum/internal/common/servs"
	"forum/internal/common/system/projfs"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	logger := slog.New(slog.NewTextHandler(
		os.Stderr,
		&slog.HandlerOptions{
			AddSource: false,
			Level:     slog.LevelInfo,
		}),
	)

	slog.SetDefault(logger)

	db, err := InitSQLite3DB()
	if err != nil {
		slog.Error(err.Error())
		return
	}
	userRepo := storage.NewUserRepo(db)

	regCommonHandlers()
	regIAMHandlers(userRepo)

	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		slog.Error(err.Error())
	}
}

func regCommonHandlers() {
	http.HandleFunc("/", handlers.NotFoundOrRedirect)

	staticDir := http.Dir(projfs.StaticDir())
	http.Handle("GET /static/",
		http.StripPrefix("/static/", http.FileServer(staticDir)))
}

func regIAMHandlers(userRepo *storage.UserRepo) {
	iamHandlers := handlers.IAMHandlers{
		IAMService: servs.NewIAMService(
			userRepo,
			passhash.NewPasswordHasher(bcrypt.DefaultCost)),
	}

	http.HandleFunc("GET /register", iamHandlers.RegisterPage)
	http.HandleFunc("POST /register", iamHandlers.Register)

	http.HandleFunc("GET /login", iamHandlers.LogInPage)
	http.HandleFunc("POST /login", iamHandlers.LogIn)
}

func InitSQLite3DB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", string(projfs.StorageDir().Join("database.db")))
	if err != nil {
		return nil, err
	}

	query := ``

	_, err = db.Exec(query)
	return db, err
}
