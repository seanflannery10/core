package helpers

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"runtime/debug"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/seanflannery10/core/internal/data"
	"github.com/seanflannery10/core/pkg/errs"
	"github.com/seanflannery10/core/pkg/validator"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	oteltrace "go.opentelemetry.io/otel/trace"
	"golang.org/x/exp/slog"
)

var ErrInvalidIDParameter = errors.New("invalid id parameter")

type contextKey string

const (
	UserContextKey = contextKey("user")
)

func ContextGetUser(r *http.Request) data.User {
	u, ok := r.Context().Value(UserContextKey).(data.User)
	if !ok {
		panic("missing user value in request context")
	}

	return u
}

func CheckAndBind(w http.ResponseWriter, r *http.Request, b render.Binder) bool {
	err := render.Bind(r, b)
	if err != nil {
		switch {
		case errors.As(err, &validator.ErrValidation):
			_ = render.Render(w, r, errs.ErrFailedValidation(validator.ErrValidation.Errors))
		case errors.Is(err, ErrInvalidIDParameter):
			_ = render.Render(w, r, errs.ErrNotFound)
		default:
			_ = render.Render(w, r, errs.ErrBadRequest(err))
		}

		return true
	}

	return false
}

func RenderAndCheck(w http.ResponseWriter, r *http.Request, ren render.Renderer) {
	span := oteltrace.SpanFromContext(r.Context())

	status, ok := r.Context().Value(render.StatusCtxKey).(int)
	if !ok {
		panic("missing status value in request context")
	}

	span.SetAttributes(semconv.HTTPStatusCodeKey.Int(status))

	spanStatus, spanMessage := semconv.SpanStatusFromHTTPStatusCode(status)
	span.SetStatus(spanStatus, spanMessage)

	err := render.Render(w, r, ren)
	if err != nil {
		slog.Error("render error", err)
	}

	span.End()
}

func ErrFuncWrapper(er *errs.ErrResponse) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		_ = render.Render(w, r, er)
	}
}

func ReadIDParam(r *http.Request) (int64, error) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil || id < 1 {
		return 0, ErrInvalidIDParameter
	}

	return id, nil
}

func ReadStringParam(qs url.Values, key string, defaultValue string) string {
	s := qs.Get(key)

	if s == "" {
		return defaultValue
	}

	return s
}

func ReadIntParam(qs url.Values, key string, defaultValue int, v *validator.Validator) int {
	s := qs.Get(key)

	if s == "" {
		return defaultValue
	}

	i, err := strconv.Atoi(s)
	if err != nil {
		v.AddError(key, "must be an integer value")
		return defaultValue
	}

	return i
}

func GetVersion() string {
	var (
		revision string
		modified bool
	)

	bi, ok := debug.ReadBuildInfo()
	if ok {
		for _, s := range bi.Settings {
			switch s.Key {
			case "vcs.revision":
				revision = s.Value
			case "vcs.modified":
				if s.Value == "true" {
					modified = true
				}
			}
		}
	}

	if revision == "" {
		return "unavailable"
	}

	if modified {
		return fmt.Sprintf("%s-dirty", revision)
	}

	return revision
}
