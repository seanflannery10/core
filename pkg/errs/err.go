package errs

import (
	"net/http"

	"github.com/go-chi/render"
	"golang.org/x/exp/slog"
)

type ErrResponse struct {
	Err       error  `json:"-"`
	Code      int    `json:"-"`
	Status    string `json:"status"`
	ErrorText string `json:"error,omitempty"`
}

func (e ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.Code)
	return nil
}

var ErrNotFound = ErrResponse{
	Code:   http.StatusNotFound,
	Status: "the requested resource could not be found",
}

var ErrMethodNotAllowed = ErrResponse{
	Code:   http.StatusMethodNotAllowed,
	Status: "the used method is not supported for this resource",
}

var ErrEditConflict = ErrResponse{
	Code:   http.StatusConflict,
	Status: "unable to update the record due to an edit conflict, please try again",
}

var ErrInvalidCredentials = ErrResponse{
	Code:   http.StatusUnauthorized,
	Status: "invalid authentication credentials",
}

var ErrInvalidAuthenticationToken = ErrResponse{
	Code:   http.StatusUnauthorized,
	Status: "invalid or missing authentication token",
}

var ErrAuthenticationRequired = ErrResponse{
	Code:   http.StatusUnauthorized,
	Status: "you must be authenticated to access this resource",
}

var ErrUserExists = ErrResponse{
	Code:   http.StatusUnprocessableEntity,
	Status: "a user with this email address already exists",
}

func ErrFailedValidation(err error) render.Renderer {
	return ErrResponse{
		Err:       err,
		Code:      http.StatusUnprocessableEntity,
		Status:    "validation failed",
		ErrorText: err.Error(),
	}
}

func ErrBadRequest(err error) render.Renderer {
	return ErrResponse{
		Err:       err,
		Code:      http.StatusBadRequest,
		Status:    "bad request",
		ErrorText: err.Error(),
	}
}

func ErrServerError(err error) render.Renderer {
	slog.Error("server error", err)

	return ErrResponse{
		Code:   http.StatusInternalServerError,
		Status: "the server encountered a problem and could not process your json",
	}
}
