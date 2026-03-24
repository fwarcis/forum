package cookies

import (
	"net/http"
	"time"
)

func Cookie(r *http.Request, name string) (*http.Cookie, error) {
	cookie, errNoCookie := r.Cookie(name)
	if errNoCookie != nil {
		return nil, errNoCookie
	}
	errInvalid := cookie.Valid()
	if errInvalid != nil {
		return nil, errInvalid
	}
	return cookie, nil
}

func NewCookie(
	name string,
	value string,
	expires time.Time,
) *http.Cookie {
	return &http.Cookie{
		Name:     name,
		Value:    value,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   int(time.Until(expires).Seconds()),
		Expires:  expires,
	}
}
