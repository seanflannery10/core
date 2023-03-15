package errs

import (
	"fmt"
	"net/http"

	"github.com/go-chi/render"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	oteltrace "go.opentelemetry.io/otel/trace"
	"golang.org/x/exp/slog"
)

type ErrResponse struct {
	Err            error       `json:"-"`
	Headers        http.Header `json:"-"`
	AppCode        codes.Code  `json:"-"`
	HTTPStatusCode int         `json:"-"`

	Message         string            `json:"message"`
	ErrorText       string            `json:"error,omitempty"`
	ValidatorErrors map[string]string `json:"errors,omitempty"`
}

func (err ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	if len(err.Headers) != 0 {
		for key, value := range err.Headers {
			w.Header()[key] = value
		}
	}

	render.Status(r, err.HTTPStatusCode)

	span := oteltrace.SpanFromContext(r.Context())
	span.SetAttributes(semconv.HTTPStatusCodeKey.Int(err.HTTPStatusCode))
	span.SetStatus(err.AppCode, "")

	span.SetAttributes(attribute.String("error.message", err.Message))

	if len(err.ErrorText) != 0 {
		span.SetAttributes(attribute.String("error.text", err.ErrorText))
	}

	if len(err.ValidatorErrors) != 0 {
		for k, v := range err.ValidatorErrors {
			k = fmt.Sprintf("error.validation.%s", k)
			span.SetAttributes(attribute.String(k, v))
		}
	}

	span.End()

	return nil
}

var ErrNotFound = &ErrResponse{ //nolint:gochecknoglobals
	AppCode:        0,
	HTTPStatusCode: http.StatusNotFound,
	Message:        "the requested resource could not be found",
}

var ErrMethodNotAllowed = &ErrResponse{ //nolint:gochecknoglobals
	AppCode:        0,
	HTTPStatusCode: http.StatusMethodNotAllowed,
	Message:        "the used method is not supported for this resource",
}

var ErrEditConflict = &ErrResponse{ //nolint:gochecknoglobals
	AppCode:        0,
	HTTPStatusCode: http.StatusConflict,
	Message:        "unable to update the record due to an edit conflict, please try again",
}

var ErrCookieNotFound = &ErrResponse{ //nolint:gochecknoglobals
	AppCode:        0,
	HTTPStatusCode: http.StatusBadRequest,
	Message:        "cookie not found",
}

var ErrInvalidCookie = &ErrResponse{ //nolint:gochecknoglobals
	AppCode:        0,
	HTTPStatusCode: http.StatusBadRequest,
	Message:        "invalid cookie",
}

var ErrInvalidCredentials = &ErrResponse{ //nolint:gochecknoglobals
	AppCode:        0,
	HTTPStatusCode: http.StatusUnauthorized,
	Message:        "invalid authentication credentials",
}

var ErrAuthenticationRequired = &ErrResponse{ //nolint:gochecknoglobals
	AppCode:        0,
	HTTPStatusCode: http.StatusUnauthorized,
	Message:        "you must be authenticated to access this resource",
}

var ErrReusedRefreshToken = &ErrResponse{ //nolint:gochecknoglobals
	AppCode:        1,
	HTTPStatusCode: http.StatusUnauthorized,
	Message:        "invalid or missing refresh token",
}

func ErrInvalidAccessToken() render.Renderer {
	headers := make(http.Header)
	headers.Set("WWW-Authenticate", "Bearer")

	return &ErrResponse{
		AppCode:        0,
		HTTPStatusCode: http.StatusUnauthorized,
		Message:        "invalid or missing access token",
		Headers:        headers,
	}
}

func ErrFailedValidation(validatorErrors map[string]string) render.Renderer {
	return &ErrResponse{
		AppCode:         0,
		HTTPStatusCode:  http.StatusUnprocessableEntity,
		Message:         "validation failed",
		ValidatorErrors: validatorErrors,
	}
}

func ErrBadRequest(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		AppCode:        0,
		HTTPStatusCode: http.StatusBadRequest,
		Message:        "bad request",
		ErrorText:      err.Error(),
	}
}

func ErrServerError(err error) render.Renderer {
	slog.Error("server error", err)

	return &ErrResponse{
		AppCode:        1,
		HTTPStatusCode: http.StatusInternalServerError,
		Message:        "the server encountered a problem and could not process your json",
	}
}
