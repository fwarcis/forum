package secur

import (
	"net/http"
	"time"
)

func NewCookie(
	name string,
	value string,
	expires time.Time,
) http.Cookie {
	return http.Cookie{
		Name:     name,
		Value:    value,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   int(time.Until(expires).Seconds()),
		Expires:  expires,
	}
}
