package handlers

import (
	"log/slog"
	"net/http"

	"forum/internal/servs"
	"forum/internal/website/security"
	"forum/internal/website/templs"
)

func UserRegisterPage(w http.ResponseWriter, r *http.Request) {
	err := templs.ExecuteWithLayout(w, "register", nil)
	if err != nil {
		slog.Error(err.Error())
	}
}

func UserLoginPage(w http.ResponseWriter, r *http.Request) {
	err := templs.ExecuteWithLayout(w, "login", nil)
	if err != nil {
		slog.Error(err.Error())
		return
	}
}

func UserRegister(srv *servs.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != "application/x-www-form-urlencoded" {
			badRequestError(w, "Wrong Content-Type")
			return
		}

		session, err := srv.Register(
			r.Context(),
			r.FormValue("login"),
			r.FormValue("email"),
			[]byte(r.FormValue("password")))
		if err != nil {
			internalServerError(w, err)
			return
		}
		sessnCookie := newSessionCookie(*session)
		http.SetCookie(w, &sessnCookie)

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func UserLogIn(srv *servs.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != "application/x-www-form-urlencoded" {
			badRequestError(w, "Wrong Content-Type")
			return
		}

		session, err := srv.LogIn(
			r.Context(),
			r.FormValue("login"),
			[]byte(r.FormValue("password")))
		if err != nil {
			internalServerError(w, err)
			return
		}
		sessnCookie := newSessionCookie(*session)
		http.SetCookie(w, &sessnCookie)

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func newSessionCookie(sessn servs.SessionData) http.Cookie {
	return security.NewCookie("session_id", sessn.ID, sessn.ExpiresAt)
}
