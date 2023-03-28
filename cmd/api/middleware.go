package main

import (
	"net/http"
	"runtime/debug"

	"github.com/go-faster/errors"
	"github.com/ogen-go/ogen/middleware"
	"golang.org/x/exp/slog"
)

var errServerError = errors.New("the server encountered a problem and could not process your request")

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
			return middleware.Response{}, errServerError
		}

		return next(req)
	}
}
