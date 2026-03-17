package handlers

import (
	"log/slog"
	"net/http"

	"forum/internal/website/templs"
)

func internalServerErrorPlainMsg(w http.ResponseWriter, msg string) {
	http.Error(w, msg, http.StatusInternalServerError)
}

func badRequestError(w http.ResponseWriter, details string) error {
	w.WriteHeader(http.StatusBadRequest)
	return templs.ExecuteError(w, templs.ErrorData{
		Message: "400 Bad Request",
		Details: details,
	})
}

func internalServerError(w http.ResponseWriter, err error) {
	slog.Error(err.Error())
	w.WriteHeader(http.StatusInternalServerError)
	err = templs.ExecuteError(w, templs.ErrorData{
		Message: "500 Internal Server Error",
		Details: err.Error(),
	})
	if err != nil {
		slog.Error(err.Error())
		internalServerErrorPlainMsg(w, err.Error())
	}
}

func NotFoundOrRedirect(w http.ResponseWriter, r *http.Request) {
}
