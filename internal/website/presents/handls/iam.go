package handls

import (
	"log/slog"
	"net/http"

	"forum/internal/common/domain/servs"
	"forum/internal/website/presents/cookies"
	"forum/internal/website/presents/errwriter"
	"forum/internal/website/presents/templs"
)

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

func (h IAMHandlers) Register(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/x-www-form-urlencoded" {
		errwriter.BadRequestError(w, "Wrong Content-Type")
		return
	}

	session, err := h.IAMService.Register(
		r.Context(),
		r.FormValue("login"),
		r.FormValue("email"),
		r.FormValue("password"))
	if err != nil {
		errwriter.InternalError(w, err)
		return
	}
	http.SetCookie(w, h.newSessionIDCookie(*session))

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h IAMHandlers) LogIn(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/x-www-form-urlencoded" {
		errwriter.BadRequestError(w, "Wrong Content-Type")
		return
	}

	session, err := h.IAMService.LogIn(
		r.Context(),
		r.FormValue("login"),
		r.FormValue("password"))
	if err != nil {
		errwriter.InternalError(w, err)
		return
	}
	http.SetCookie(w, h.newSessionIDCookie(*session))

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h IAMHandlers) newSessionIDCookie(sess servs.SessionData) *http.Cookie {
	return cookies.NewCookie("session_id", sess.ID, sess.ExpiresAt)
}

func (h IAMHandlers) sessionIDCookie(r *http.Request) (*http.Cookie, error) {
	return cookies.Cookie(r, "session_id")
}
