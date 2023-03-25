package middleware

import (
	"crypto/sha256"
	"net/http"
	"runtime/debug"
	"strings"
	"time"

	"golang.org/x/exp/slog"

	"github.com/ogen-go/ogen/middleware"
	"github.com/seanflannery10/core/internal/data"

	"github.com/go-faster/errors"
	"github.com/seanflannery10/core/internal/shared/utils"
)

var errInvalidAccesToken = errors.New("invalid access token")

func Authenticate(q data.Queries) middleware.Middleware {
	return func(req middleware.Request, next func(req middleware.Request) (middleware.Response, error)) (middleware.Response, error) {
		req.Raw.Header.Add("Vary", "Authorization")

		authorizationHeader := req.Raw.Header.Get("Authorization")

		switch authorizationHeader {
		case "":
			user := data.AnonymousUser
			utils.ContextSetUser(&req, user)
		default:
			headerParts := strings.Split(authorizationHeader, " ")
			if len(headerParts) != 2 || headerParts[0] != "Bearer" { //nolint:revive
				return middleware.Response{}, errInvalidAccesToken
			}

			token := headerParts[1] //nolint:revive
			tokenHash := sha256.Sum256([]byte(token))

			user, err := q.GetUserFromToken(req.Context, data.GetUserFromTokenParams{
				Hash:   tokenHash[:],
				Scope:  data.ScopeAccess,
				Expiry: time.Now(),
			})
			if err != nil {
				return middleware.Response{}, errInvalidAccesToken
			}

			utils.ContextSetUser(&req, &user)
		}

		resp, err := next(req)
		return resp, err
	}
}

func RecoverPanic() middleware.Middleware {
	return func(req middleware.Request, next func(req middleware.Request) (middleware.Response, error)) (middleware.Response, error) {
		defer func() {
			if rvr := recover(); rvr != nil {
				if rvr == http.ErrAbortHandler { //nolint:errorlint,goerr113
					panic(rvr)
				}

				slog.Log(
					req.Context,
					slog.LevelError,
					"panic recovery error",
					"error",
					rvr,
					"stack",
					string(debug.Stack()),
				)

				if req.Raw.Header.Get("Connection") != "Upgrade" {
					// TODO Return 500 error
				}
			}
		}()

		resp, err := next(req)
		return resp, err
	}
}

func RequireAuthenticatedUser() middleware.Middleware {
	return func(req middleware.Request, next func(req middleware.Request) (middleware.Response, error)) (middleware.Response, error) {
		user := utils.ContextGetUser(&req)

		if user.IsAnonymous() {
			// TODO return authenticationRequired error
		}

		resp, err := next(req)
		return resp, err
	}
}

func RequireActivatedUser() middleware.Middleware {
	return func(req middleware.Request, next func(req middleware.Request) (middleware.Response, error)) (middleware.Response, error) {
		user := utils.ContextGetUser(&req)

		if !user.Activated {
			// TODO return actiaved user required
		}

		resp, err := next(req)
		return resp, err
	}
}
