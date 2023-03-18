package middleware

import (
	"net/http"
	"runtime/debug"

	"github.com/go-chi/render"
	"github.com/seanflannery10/core/internal/pkg/errs"
	"golang.org/x/exp/slog"
)

func RecoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rvr := recover(); rvr != nil {
				if rvr == http.ErrAbortHandler { //nolint:errorlint,goerr113
					panic(rvr)
				}

				slog.Log(
					r.Context(),
					slog.LevelError,
					"panic recovery error",
					"error",
					rvr,
					"stack",
					string(debug.Stack()),
				)

				if r.Header.Get("Connection") != "Upgrade" {
					w.WriteHeader(http.StatusInternalServerError)
				}

				render.JSON(w, r, &errs.ErrResponse{
					Message: "the server encountered a problem and could not process your json",
				})
			}
		}()

		next.ServeHTTP(w, r)
	})
}
