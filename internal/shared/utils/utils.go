package utils

import (
	"context"
	"fmt"
	"runtime/debug"

	"github.com/seanflannery10/core/internal/generated/data"
)

const (
	userContextKey = contextKey("user")
)

type contextKey string

func ContextSetUser(ctx context.Context, user *data.User) context.Context {
	return context.WithValue(ctx, userContextKey, *user)
}

func ContextSetCookieValue(ctx context.Context, s string) context.Context {
	return context.WithValue(ctx, userContextKey, s)
}

func ContextGetUser(ctx context.Context) data.User {
	user, ok := ctx.Value(userContextKey).(data.User)
	if !ok {
		panic("missing user value in request context")
	}

	return user
}

func ContextGetCookieValue(ctx context.Context) string {
	cookieValue, ok := ctx.Value(userContextKey).(string)
	if !ok {
		panic("missing cookie value in request context")
	}

	return cookieValue
}

func GetVersion() string {
	var (
		revision string
		modified bool
	)

	bi, ok := debug.ReadBuildInfo()
	if ok {
		for _, s := range bi.Settings {
			switch s.Key {
			case "vcs.revision":
				revision = s.Value
			case "vcs.modified":
				if s.Value == "true" {
					modified = true
				}
			}
		}
	}

	if revision == "" {
		return "unavailable"
	}

	if modified {
		return fmt.Sprintf("%s-dirty", revision)
	}

	return revision
}
