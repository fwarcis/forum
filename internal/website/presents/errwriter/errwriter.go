package errwriter

import (
	"log/slog"
	"net/http"

	"forum/internal/website/presents/templs"
)

func BadRequestError(w http.ResponseWriter, details string) error {
	w.WriteHeader(http.StatusBadRequest)
	return templs.ExecuteError(w, templs.ErrorData{
		Message: "400 Bad Request",
		Details: details,
	})
}

func InternalError(w http.ResponseWriter, err error) {
	slog.Error(err.Error())
	w.WriteHeader(http.StatusInternalServerError)
	err = templs.ExecuteError(w, templs.ErrorData{
		Message: "500 Internal Server Error",
		Details: err.Error(),
	})
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func NotFoundOrRedirect(w http.ResponseWriter, r *http.Request) {
}
