package main

import (
	"log/slog"
	"net/http"
	"os"

	"forum/internal/projfs"
	"forum/internal/servs"
	"forum/internal/storage"
	"forum/internal/website/handlers"

	_ "github.com/mattn/go-sqlite3"
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

	db, err := storage.InitSQLite3DB()
	if err != nil {
		slog.Error(err.Error())
		return
	}

	userRepo := storage.NewSQLite3UserRepo(db)

	http.HandleFunc("/", handlers.NotFoundOrRedirect)

	staticDir := http.Dir(projfs.StaticDir())
	http.Handle("GET /static/",
		http.StripPrefix("/static/", http.FileServer(staticDir)))

	http.HandleFunc("GET /register", handlers.UserRegisterPage)
	http.HandleFunc("POST /register",
		handlers.UserRegister(servs.NewUserService(userRepo)))

	http.HandleFunc("GET /login", handlers.UserLoginPage)
	http.HandleFunc("POST /login",
		handlers.UserLogIn(servs.NewUserService(userRepo)))

	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		slog.Error(err.Error())
	}
}
