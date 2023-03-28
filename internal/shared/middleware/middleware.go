package middleware

import (
	"context"
	"crypto/sha256"
	"net/http"
	"regexp"
	"runtime/debug"
	"strings"
	"time"

	"github.com/go-faster/errors"
	"github.com/go-faster/jx"
	"github.com/ogen-go/ogen/middleware"
	"github.com/ogen-go/ogen/ogenerrors"
	"github.com/seanflannery10/core/internal/api"
	"github.com/seanflannery10/core/internal/data"
	"github.com/seanflannery10/core/internal/shared/utils"
	"golang.org/x/exp/slog"
)

var (
	errInvalidAccessToken   = errors.New("invalid or missing authentication token")
	errServerError          = errors.New("the server encountered a problem and could not process your request")
	errUserNotActivated     = errors.New("your user account must be activated to access this resource")
	errUserNotAuthenticated = errors.New("you must be authenticated to access this resource")
)

type Security struct {
	Queries *data.Queries
}

func (s *Security) HandleAccess(ctx context.Context, _ string, t api.Access) (context.Context, error) {
	switch t.Token {
	case "":
		return utils.ContextSetUser(ctx, data.AnonymousUser), nil
	default:
		tokenHash := sha256.Sum256([]byte(t.Token))

		user, err := s.Queries.GetUserFromToken(ctx, data.GetUserFromTokenParams{
			Hash:   tokenHash[:],
			Scope:  data.ScopeAccess,
			Expiry: time.Now(),
		})
		if err != nil {
			return ctx, errInvalidAccessToken
		}

		return utils.ContextSetUser(ctx, &user), nil
	}
}

func RecoverPanic() middleware.Middleware {
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

func RequireAuthenticatedUser() middleware.Middleware {
	return func(req middleware.Request, next func(req middleware.Request) (middleware.Response, error)) (middleware.Response, error) {
		user := utils.ContextGetUser(req.Context)

		if user.IsAnonymous() {
			return middleware.Response{}, errUserNotAuthenticated
		}

		return next(req)
	}
}

func RequireActivatedUser() middleware.Middleware {
	return func(req middleware.Request, next func(req middleware.Request) (middleware.Response, error)) (middleware.Response, error) {
		user := utils.ContextGetUser(req.Context)

		if !user.Activated {
			return middleware.Response{}, errUserNotActivated
		}

		return next(req)
	}
}

func ErrorHandler(_ context.Context, w http.ResponseWriter, _ *http.Request, err error) {
	code := ogenerrors.ErrorCode(err)

	switch {
	case errors.Is(err, errInvalidAccessToken):
		w.Header().Set("WWW-Authenticate", "Bearer")
		code = http.StatusUnauthorized
	case errors.Is(err, errUserNotAuthenticated):
		code = http.StatusUnauthorized
	case errors.Is(err, errUserNotActivated):
		code = http.StatusForbidden
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	errStr := strings.ReplaceAll(err.Error(), "\"", "")

	re := regexp.MustCompile(`^operation \w*: (.+)$`)
	errSubmatch := re.FindStringSubmatch(errStr)

	e := jx.GetEncoder()
	e.ObjStart()
	e.FieldStart("message")
	e.StrEscape("request error")
	e.FieldStart("error")
	e.StrEscape(errSubmatch[1])
	e.ObjEnd()

	_, _ = w.Write(e.Bytes())
}
