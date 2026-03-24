package middlews

import (
	"context"
	"net/http"

	"forum/internal/common/domain/servs"
	"forum/internal/website/presents/templs"

	"github.com/google/uuid"
)

type key string

var userKey key

type IAMMiddlewares struct {
	IAMService      *servs.IAMService
	strToUUIDTransr string
	next            http.HandlerFunc
}

func (m IAMMiddlewares) Auth(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err == nil && cookie.Valid() == nil {
		sessionID, err := uuid.Parse(cookie.Value)
		var ctx context.Context
		if err == nil {
			ctx, err = m.IAMService.Auth(r.Context(), sessionID)
		}
		if err == nil {
			m.next(w, r.WithContext(ctx))
			return
		}
	}
	templs.ExecuteWithLayout(w, "login", nil)
}
