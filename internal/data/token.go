package data

import (
	"context"
	"crypto/sha256"
	"fmt"
	"time"
)

const (
	ScopeAccess        = "access"
	ScopeActivation    = "activation"
	ScopePasswordReset = "password-reset"
	ScopeRefresh       = "refresh"
)

func (q *Queries) GetUserFromTokenHelper(ctx context.Context, tokenPlaintext, scope string) (User, error) {
	tokenHash := sha256.Sum256([]byte(tokenPlaintext))

	user, err := q.GetUserFromToken(ctx, GetUserFromTokenParams{
		Hash:   tokenHash[:],
		Scope:  scope,
		Expiry: time.Now(),
	})
	if err != nil {
		return User{}, fmt.Errorf("failed get user from token: %w", err)
	}

	return user, nil
}
