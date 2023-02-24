package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/seanflannery10/core/pkg/assert"
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

// func TestMiddleware_RequireAuthenticatedUser(t *testing.T) {
//	t.Run("Bad Auth", func(t *testing.T) {
//		rr := httptest.NewRecorder()
//		r, _ := http.NewRequestWithContext(ctx, http.MethodGet, "/", nil)
//
//		RequireAuthenticatedUser(next).ServeHTTP(rr, r)
//
//		assert.Contains(t, rr.Body.String(), "you must be authenticated to access this resource")
//		assert.Equal(t, rr.Code, http.StatusUnauthorized)
//	})
//
//	t.Run("Good Auth", func(t *testing.T) {
//		rr := httptest.NewRecorder()
//		r, _ := http.NewRequestWithContext(ctx, http.MethodGet, "/", nil)
//
//		r = r.WithContext(context.WithValue(r.Context(), ctxkeys.UserContextKey, "Test"))
//
//		RequireAuthenticatedUser(next).ServeHTTP(rr, r)
//
//		assert.Equal(t, rr.Body.String(), "OK")
//		assert.Equal(t, rr.Code, http.StatusOK)
//	})
//}

func TestMiddleware_Metrics(t *testing.T) {
	rr := httptest.NewRecorder()
	r, _ := http.NewRequestWithContext(ctx, http.MethodGet, "/", nil)

	Metrics(next).ServeHTTP(rr, r)

	assert.Equal(t, rr.Body.String(), "OK")
}

func TestMiddleware_RecoverPanic(t *testing.T) {
	t.Run("No Panic", func(t *testing.T) {
		rr := httptest.NewRecorder()
		r, _ := http.NewRequestWithContext(ctx, http.MethodGet, "/", nil)

		RecoverPanic(next).ServeHTTP(rr, r)

		assert.Equal(t, rr.Body.String(), "OK")
	})

	t.Run("Panic", func(t *testing.T) {
		rr := httptest.NewRecorder()
		r, _ := http.NewRequestWithContext(ctx, http.MethodGet, "/", nil)

		homeHandler := http.HandlerFunc(func(http.ResponseWriter, *http.Request) { panic("test error") })

		RecoverPanic(homeHandler).ServeHTTP(rr, r)

		assert.Contains(t, rr.Body.String(), "the server encountered a problem and could not process your json")
	})
}
