package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/seanflannery10/core/internal/assert"
)

var ctx = context.Background()

func TestHealthcheck(t *testing.T) {
	rr := httptest.NewRecorder()

	r, err := http.NewRequestWithContext(ctx, http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	Healthcheck(rr, r)

	res := rr.Result()
	defer res.Body.Close()

	assert.Equal(t, res.StatusCode, http.StatusOK)
	assert.Contains(t, rr.Body.String(), "available")
}
