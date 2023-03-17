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

const (
	appNormal = 0
	appError  = 1
	noLength  = 0
)

type ErrResponse struct { //nolint:govet
	ValidatorErrors map[string]string `json:"errors,omitempty"`
	Message         string            `json:"message"`
	ErrorText       string            `json:"error,omitempty"`

	Err            error       `json:"-"`
	Headers        http.Header `json:"-"`
	AppCode        codes.Code  `json:"-"`
	HTTPStatusCode int         `json:"-"`
}

func (err *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	if len(err.Headers) != noLength {
		for key, value := range err.Headers {
			w.Header()[key] = value
		}
	}

	render.Status(r, err.HTTPStatusCode)

	span := oteltrace.SpanFromContext(r.Context())
	span.SetAttributes(semconv.HTTPStatusCodeKey.Int(err.HTTPStatusCode))
	span.SetStatus(err.AppCode, "")

	span.SetAttributes(attribute.String("error.message", err.Message))

	if err.ErrorText != "" {
		span.SetAttributes(attribute.String("error.text", err.ErrorText))
	}

	if len(err.ValidatorErrors) != appNormal {
		for k, v := range err.ValidatorErrors {
			k = fmt.Sprintf("error.validation.%s", k)
			span.SetAttributes(attribute.String(k, v))
		}
	}

	span.End()

	return nil
}

func ErrNotFound() render.Renderer {
	return &ErrResponse{
		AppCode:        appNormal,
		HTTPStatusCode: http.StatusNotFound,
		Message:        "the requested resource could not be found",
	}
}

func ErrMethodNotAllowed() render.Renderer {
	return &ErrResponse{
		AppCode:        appNormal,
		HTTPStatusCode: http.StatusMethodNotAllowed,
		Message:        "the used method is not supported for this resource",
	}
}

func ErrEditConflict() render.Renderer {
	return &ErrResponse{
		AppCode:        appNormal,
		HTTPStatusCode: http.StatusConflict,
		Message:        "unable to update the record due to an edit conflict, please try again",
	}
}

func ErrCookieNotFound() render.Renderer {
	return &ErrResponse{
		AppCode:        appNormal,
		HTTPStatusCode: http.StatusBadRequest,
		Message:        "cookie not found",
	}
}

func ErrInvalidCookie() render.Renderer {
	return &ErrResponse{
		AppCode:        appNormal,
		HTTPStatusCode: http.StatusBadRequest,
		Message:        "invalid cookie",
	}
}

func ErrInvalidCredentials() render.Renderer {
	return &ErrResponse{
		AppCode:        appNormal,
		HTTPStatusCode: http.StatusUnauthorized,
		Message:        "invalid authentication credentials",
	}
}

func ErrAuthenticationRequired() render.Renderer {
	return &ErrResponse{
		AppCode:        appNormal,
		HTTPStatusCode: http.StatusUnauthorized,
		Message:        "you must be authenticated to access this resource",
	}
}

func ErrReusedRefreshToken() render.Renderer {
	return &ErrResponse{
		AppCode:        appError,
		HTTPStatusCode: http.StatusUnauthorized,
		Message:        "invalid or missing refresh token",
	}
}

func ErrInvalidAccessToken() render.Renderer {
	headers := make(http.Header)
	headers.Set("WWW-Authenticate", "Bearer")

	return &ErrResponse{
		AppCode:        appNormal,
		HTTPStatusCode: http.StatusUnauthorized,
		Message:        "invalid or missing access token",
		Headers:        headers,
	}
}

func ErrFailedValidation(validatorErrors map[string]string) render.Renderer {
	return &ErrResponse{
		AppCode:         appNormal,
		HTTPStatusCode:  http.StatusUnprocessableEntity,
		Message:         "validation failed",
		ValidatorErrors: validatorErrors,
	}
}

func ErrBadRequest(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		AppCode:        appNormal,
		HTTPStatusCode: http.StatusBadRequest,
		Message:        "bad request",
		ErrorText:      err.Error(),
	}
}

func ErrServerError(err error) render.Renderer {
	slog.Error("server error", err)

	return &ErrResponse{
		AppCode:        appError,
		HTTPStatusCode: http.StatusInternalServerError,
		Message:        "the server encountered a problem and could not process your json",
	}
}
