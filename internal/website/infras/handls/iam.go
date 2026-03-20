package handls

import (
	"log/slog"
	"net/http"

	"forum/internal/common/servs"
	"forum/internal/website/infras/secur"
	"forum/internal/website/views/templs"
)

const sessionCookieName = "session_id"

func newSessionCookie(sessn servs.SessionData) http.Cookie {
	return secur.NewCookie(sessionCookieName, sessn.ID, sessn.ExpiresAt)
}

type IAMHandlers struct {
	IAMService *servs.IAMService
}

func (h IAMHandlers) RegisterPage(w http.ResponseWriter, r *http.Request) {
	err := templs.ExecuteWithLayout(w, "register", nil)
	if err != nil {
		slog.Error(err.Error())
	}
}

func (h IAMHandlers) LogInPage(w http.ResponseWriter, r *http.Request) {
	err := templs.ExecuteWithLayout(w, "login", nil)
	if err != nil {
		slog.Error(err.Error())
		return
	}
}

func (h IAMHandlers) CheckSession(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(sessionCookieName)
	if err == nil && cookie.Valid() == nil {
		return
	}
}

func (h IAMHandlers) Register(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/x-www-form-urlencoded" {
		badRequestError(w, "Wrong Content-Type")
		return
	}

	session, err := h.IAMService.Register(
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

func (h IAMHandlers) LogIn(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/x-www-form-urlencoded" {
		badRequestError(w, "Wrong Content-Type")
		return
	}

	session, err := h.IAMService.LogIn(
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
