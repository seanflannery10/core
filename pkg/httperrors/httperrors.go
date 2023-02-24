package httperrors

import (
	"fmt"
	"net/http"

	"github.com/seanflannery10/core/pkg/helpers"
	"github.com/seanflannery10/core/pkg/validator"
	"golang.org/x/exp/slog"
)

func ErrorMessage(w http.ResponseWriter, r *http.Request, status int, message string) {
	ErrorMessageWithHeaders(w, r, status, message, nil)
}

func ErrorMessageWithHeaders(w http.ResponseWriter, r *http.Request, status int, message string, headers http.Header) {
	err := helpers.WriteJSONWithHeaders(w, status, map[string]string{"error": message}, headers)
	if err != nil {
		slog.Error(
			"json write error",
			err,
			"request_method",
			r.Method,
			"request_url",
			r.URL.String(),
		)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func FailedValidation(w http.ResponseWriter, r *http.Request, v *validator.Validator) {
	err := helpers.WriteJSON(w, http.StatusUnprocessableEntity, map[string]map[string]string{"error": v.Errors})
	if err != nil {
		slog.Error(
			"json write error",
			err,
			"request_method",
			r.Method,
			"request_url",
			r.URL.String(),
		)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func ServerError(w http.ResponseWriter, r *http.Request, err error) {
	slog.Error("server error", err)

	message := "the server encountered a problem and could not process your json"
	ErrorMessage(w, r, http.StatusInternalServerError, message)
}

func NotFound(w http.ResponseWriter, r *http.Request) {
	message := "the requested resource could not be found"
	ErrorMessage(w, r, http.StatusNotFound, message)
}

func MethodNotAllowed(w http.ResponseWriter, r *http.Request) {
	message := fmt.Sprintf("The %s method is not supported for this resource", r.Method)
	ErrorMessage(w, r, http.StatusMethodNotAllowed, message)
}

func BadRequest(w http.ResponseWriter, r *http.Request, err error) {
	ErrorMessage(w, r, http.StatusBadRequest, err.Error())
}

func InvalidCredentials(w http.ResponseWriter, r *http.Request) {
	message := "invalid authentication credentials"
	ErrorMessage(w, r, http.StatusUnauthorized, message)
}

func InvalidAuthenticationToken(w http.ResponseWriter, r *http.Request) {
	headers := make(http.Header)
	headers.Set("WWW-Authenticate", "Bearer")

	ErrorMessageWithHeaders(w, r, http.StatusUnauthorized, "invalid or missing authentication token", headers)
}

func AuthenticationRequired(w http.ResponseWriter, r *http.Request) {
	message := "you must be authenticated to access this resource"
	ErrorMessage(w, r, http.StatusUnauthorized, message)
}

func RateLimitExceededResponse(w http.ResponseWriter, r *http.Request) {
	message := "rate limit exceeded"
	ErrorMessage(w, r, http.StatusTooManyRequests, message)
}

func EditConflict(w http.ResponseWriter, r *http.Request) {
	message := "unable to update the record due to an edit conflict, please try again"
	ErrorMessage(w, r, http.StatusConflict, message)
}
