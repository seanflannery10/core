package handler_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-faster/errors"
	"github.com/ogen-go/ogen/ogenerrors"
	"github.com/seanflannery10/core/internal/api"
	"github.com/seanflannery10/core/internal/handler"
	"github.com/seanflannery10/core/internal/logic"
	"github.com/seanflannery10/core/internal/shared/pagination"
	"github.com/stretchr/testify/assert"
)

func TestNewError(t *testing.T) {
	testCases := []struct {
		Error      error
		StatusCode int
	}{
		{Error: logic.ErrInvalidCredentials, StatusCode: http.StatusUnauthorized},
		{Error: logic.ErrReusedRefreshToken, StatusCode: http.StatusUnauthorized},
		{Error: logic.ErrEmailNotFound, StatusCode: http.StatusNotFound},
		{Error: logic.ErrMessageNotFound, StatusCode: http.StatusNotFound},
		{Error: logic.ErrEditConflict, StatusCode: http.StatusConflict},
		{Error: logic.ErrActivationRequired, StatusCode: http.StatusUnprocessableEntity},
		{Error: logic.ErrInvalidToken, StatusCode: http.StatusUnprocessableEntity},
		{Error: pagination.ErrPageValueToHigh, StatusCode: http.StatusUnprocessableEntity},
		{Error: logic.ErrUserAlreadyActivated, StatusCode: http.StatusUnprocessableEntity},
		{Error: logic.ErrUserExists, StatusCode: http.StatusUnprocessableEntity},
		{Error: logic.ErrServerError, StatusCode: http.StatusInternalServerError},
	}

	for _, tc := range testCases {
		t.Run(tc.Error.Error(), func(t *testing.T) {
			expected := &api.ErrorResponseStatusCode{
				StatusCode: tc.StatusCode,
				Response:   api.ErrorResponse{Error: tc.Error.Error()},
			}

			response := newTestHandler(t).NewError(context.Background(), errors.Wrap(tc.Error, "testing"))
			assert.Equal(t, expected, response)
		})
	}
}

func TestErrorHandler(t *testing.T) {
	testCases := []struct {
		Error      error
		StatusCode int
	}{
		{Error: logic.ErrServerError, StatusCode: http.StatusInternalServerError},
		{Error: ogenerrors.ErrSecurityRequirementIsNotSatisfied, StatusCode: http.StatusInternalServerError},
	}

	for _, tc := range testCases {
		t.Run(tc.Error.Error(), func(t *testing.T) {
			req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "/test", http.NoBody)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()

			handler.ErrorHandler(context.Background(), rr, req, errors.Wrap(tc.Error, "testing"))

			assert.Equal(t, tc.StatusCode, rr.Code)
		})
	}
}
