package main

import (
	"context"
	"crypto/sha256"
	"time"

	"github.com/go-faster/errors"
	"github.com/seanflannery10/core/internal/api"
	"github.com/seanflannery10/core/internal/data"
	"github.com/seanflannery10/core/internal/shared/utils"
)

var (
	errInvalidAccessToken = errors.New("invalid or missing authentication token")
	errUserNotActivated   = errors.New("your user account must be activated to access this resource")
)

type security struct {
	Queries *data.Queries
}

func (s *security) HandleAccess(ctx context.Context, _ string, t api.Access) (context.Context, error) {
	tokenHash := sha256.Sum256([]byte(t.Token))

	user, err := s.Queries.GetUserFromToken(ctx, data.GetUserFromTokenParams{
		Hash:   tokenHash[:],
		Scope:  data.ScopeAccess,
		Expiry: time.Now(),
	})
	if err != nil {
		return ctx, errInvalidAccessToken
	}

	if !user.Activated {
		return ctx, errUserNotActivated
	}

	return utils.ContextSetUser(ctx, &user), nil
}
