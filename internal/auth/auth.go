package auth

import (
	"context"
	"fmt"
	"net/http"
)

type key string

var userKey key

func SetUser(r *http.Request, user string) *http.Request {
	ctx := context.WithValue(r.Context(), userKey, user)
	return r.WithContext(ctx)
}

func GetUser(r *http.Request) string {
	user := r.Context().Value(userKey)
	if user == nil {
		return ""
	}

	return fmt.Sprintf("%v", user)
}
