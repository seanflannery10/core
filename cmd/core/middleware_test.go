package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/seanflannery10/core/internal/assert"
	"github.com/seanflannery10/core/internal/data"
	"github.com/seanflannery10/core/internal/helpers"
)

var (
	ctx  = context.Background()
	next = http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("OK"))
		if err != nil {
			return
		}
	}))
)

func TestMiddleware_RequireAuthenticatedUser(t *testing.T) {
	t.Run("Bad Auth", func(t *testing.T) {
		app := application{}

		rr := httptest.NewRecorder()
		r, _ := http.NewRequestWithContext(ctx, http.MethodGet, "/", nil)

		r = helpers.ContextSetUser(r, data.AnonymousUser)

		app.requireAuthenticatedUser(next).ServeHTTP(rr, r)

		assert.Contains(t, rr.Body.String(), "you must be authenticated to access this resource")
		assert.Equal(t, rr.Code, http.StatusUnauthorized)
	})

	t.Run("Good Auth", func(t *testing.T) {
		app := application{}

		rr := httptest.NewRecorder()
		r, _ := http.NewRequestWithContext(ctx, http.MethodGet, "/", nil)
		r = helpers.ContextSetUser(r, &data.User{})

		app.requireAuthenticatedUser(next).ServeHTTP(rr, r)

		assert.Equal(t, rr.Body.String(), "OK")
		assert.Equal(t, rr.Code, http.StatusOK)
	})
}

func TestMiddleware_Metrics(t *testing.T) {
	app := application{}

	rr := httptest.NewRecorder()
	r, _ := http.NewRequestWithContext(ctx, http.MethodGet, "/", nil)

	app.metrics(next).ServeHTTP(rr, r)

	assert.Equal(t, rr.Body.String(), "OK")
}

func TestMiddleware_RecoverPanic(t *testing.T) {
	t.Run("No Panic", func(t *testing.T) {
		app := application{}

		rr := httptest.NewRecorder()
		r, _ := http.NewRequestWithContext(ctx, http.MethodGet, "/", nil)

		app.recoverPanic(next).ServeHTTP(rr, r)

		assert.Equal(t, rr.Body.String(), "OK")
	})

	t.Run("Panic", func(t *testing.T) {
		app := application{}

		rr := httptest.NewRecorder()
		r, _ := http.NewRequestWithContext(ctx, http.MethodGet, "/", nil)

		homeHandler := http.HandlerFunc(func(http.ResponseWriter, *http.Request) { panic("test error") })

		app.recoverPanic(homeHandler).ServeHTTP(rr, r)

		assert.Contains(t, rr.Body.String(), "the server encountered a problem and could not process your json")
	})
}
