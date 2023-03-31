package main

import (
	"net/http"
	"runtime/debug"

	"github.com/ogen-go/ogen/middleware"
	"github.com/seanflannery10/core/internal/logic"
	"golang.org/x/exp/slog"
)

func (app *application) RecoverPanic() middleware.Middleware {
	return func(req middleware.Request, next func(req middleware.Request) (middleware.Response, error)) (middleware.Response, error) {
		recovered := false

		defer func() {
			if rvr := recover(); rvr != nil {
				if rvr == http.ErrAbortHandler { //nolint:errorlint,goerr113
					panic(rvr)
				}

				slog.Log(req.Context, slog.LevelError, "panic recovery error", "error", rvr, "stack", string(debug.Stack()))

				req.Raw.Header.Add("Connection", "close")

				recovered = true
			}
		}()

		if recovered {
			return middleware.Response{}, logic.ErrServerError
		}

		return next(req)
	}
}
