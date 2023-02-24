package httperrors

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/seanflannery10/core/pkg/assert"
	"github.com/seanflannery10/core/pkg/validator"
)

var ctx = context.Background()

func TestErrorMessage(t *testing.T) {
	tests := []struct {
		sc   int
		body string
	}{
		{200, "testing status code 200"},
		{401, "testing status code 401"},
		{404, "testing status code 404"},
		{500, "testing status code 500"},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%#v", tt.sc), func(t *testing.T) {
			rr := httptest.NewRecorder()

			r, err := http.NewRequestWithContext(ctx, http.MethodGet, "/", nil)
			if err != nil {
				t.Fatal(err)
			}

			ErrorMessage(rr, r, tt.sc, tt.body)

			res := rr.Result()
			defer res.Body.Close()

			assert.Contains(t, rr.Body.String(), tt.body)
			assert.Equal(t, res.StatusCode, tt.sc)
		})
	}
}

func TestFailedValidation(t *testing.T) {
	rr := httptest.NewRecorder()

	r, err := http.NewRequestWithContext(ctx, http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	FailedValidation(rr, r, &validator.Validator{})

	res := rr.Result()
	defer res.Body.Close()

	assert.Equal(t, res.StatusCode, http.StatusUnprocessableEntity)
}

func TestStatusCodesWithError(t *testing.T) {
	tests := []struct {
		name string
		sc   int
		f    func(http.ResponseWriter, *http.Request, error)
	}{
		{
			name: "ServerError",
			sc:   http.StatusInternalServerError,
			f:    ServerError,
		},
		{
			name: "BadRequest",
			sc:   http.StatusBadRequest,
			f:    BadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rr := httptest.NewRecorder()

			r, err := http.NewRequestWithContext(ctx, http.MethodGet, "/", nil)
			if err != nil {
				t.Fatal(err)
			}

			tt.f(rr, r, errors.New("test")) //nolint:goerr113

			res := rr.Result()
			defer res.Body.Close()

			assert.Equal(t, res.StatusCode, tt.sc)
		})
	}
}

func TestStatusCodesWithoutError(t *testing.T) {
	tests := []struct {
		name string
		sc   int
		f    func(http.ResponseWriter, *http.Request)
	}{
		{
			name: "NotFound",
			sc:   http.StatusNotFound,
			f:    NotFound,
		},
		{
			name: "MethodNotAllowed",
			sc:   http.StatusMethodNotAllowed,
			f:    MethodNotAllowed,
		},
		{
			name: "InvalidAuthenticationToken",
			sc:   http.StatusUnauthorized,
			f:    InvalidAuthenticationToken,
		},
		{
			name: "AuthenticationRequired",
			sc:   http.StatusUnauthorized,
			f:    AuthenticationRequired,
		},
		{
			name: "RateLimitExceededResponse",
			sc:   http.StatusTooManyRequests,
			f:    RateLimitExceededResponse,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rr := httptest.NewRecorder()

			r, err := http.NewRequestWithContext(ctx, http.MethodGet, "/", nil)
			if err != nil {
				t.Fatal(err)
			}

			tt.f(rr, r)

			res := rr.Result()
			defer res.Body.Close()

			assert.Equal(t, res.StatusCode, tt.sc)
		})
	}
}
