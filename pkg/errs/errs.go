package errs

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/seanflannery10/core/pkg/validator"
	"golang.org/x/exp/slog"
)

type ErrResponse struct {
	Err             error             `json:"-"`
	Code            int               `json:"-"`
	Message         string            `json:"message"`
	ErrorText       string            `json:"error,omitempty"`
	ValidatorErrors map[string]string `json:"errors,omitempty"`
}

func (err ErrResponse) Render(_ http.ResponseWriter, r *http.Request) error {
	render.Status(r, err.Code)
	return nil
}

var ErrNotFound = &ErrResponse{
	Code:    http.StatusNotFound,
	Message: "the requested resource could not be found",
}

var ErrMethodNotAllowed = &ErrResponse{
	Code:    http.StatusMethodNotAllowed,
	Message: "the used method is not supported for this resource",
}

var ErrEditConflict = &ErrResponse{
	Code:    http.StatusConflict,
	Message: "unable to update the record due to an edit conflict, please try again",
}

var ErrInvalidCredentials = &ErrResponse{
	Code:    http.StatusUnauthorized,
	Message: "invalid authentication credentials",
}

var ErrAuthenticationRequired = &ErrResponse{
	Code:    http.StatusUnauthorized,
	Message: "you must be authenticated to access this resource",
}

func ErrInvalidAuthenticationToken() render.Renderer {
	headers := make(http.Header)
	headers.Set("WWW-Authenticate", "Bearer")

	return &ErrResponse{
		Code:    http.StatusUnauthorized,
		Message: "invalid or missing authentication token",
	}
}

func ErrFailedValidation(v *validator.Validator) render.Renderer {
	return &ErrResponse{
		Code:            http.StatusUnprocessableEntity,
		Message:         "validation failed",
		ValidatorErrors: v.Errors,
	}
}

func ErrBadRequest(err error) render.Renderer {
	return &ErrResponse{
		Err:       err,
		Code:      http.StatusBadRequest,
		Message:   "bad request",
		ErrorText: err.Error(),
	}
}

func ErrServerError(err error) render.Renderer {
	slog.Error("server error", err)

	return &ErrResponse{
		Code:    http.StatusInternalServerError,
		Message: "the server encountered a problem and could not process your json",
	}
}
