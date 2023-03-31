package logic

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
	"fmt"
	"time"

	"github.com/go-faster/errors"
	"github.com/seanflannery10/core/internal/api"
	"github.com/seanflannery10/core/internal/data"
	"golang.org/x/crypto/bcrypt"
)

const (
	PasswordCost       = 13
	ScopeAccess        = "access"
	ScopeActivation    = "activation"
	ScopePasswordReset = "password-reset"
	ScopeRefresh       = "refresh"
)

func setPassword(user *data.User, plaintextPassword string) (*data.User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintextPassword), PasswordCost)
	if err != nil {
		return nil, fmt.Errorf("failed set password: %w", err)
	}

	user.PasswordHash = hash

	return user, nil
}

func comparePasswords(user *data.User, plaintextPassword string) error {
	if err := bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(plaintextPassword)); err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return ErrInvalidCredentials
		default:
			return fmt.Errorf("failed compare password: %w", err)
		}
	}

	return nil
}

func getUserFromToken(ctx context.Context, q *data.Queries, tokenPlaintext, scope string) (*data.User, error) {
	tokenHash := sha256.Sum256([]byte(tokenPlaintext))

	user, err := q.GetUserFromToken(ctx, data.GetUserFromTokenParams{Hash: tokenHash[:], Scope: scope, Expiry: time.Now()})
	if err != nil {
		return nil, fmt.Errorf("failed get user from token: %w", err)
	}

	return user, nil
}

func newToken(ctx context.Context, q *data.Queries, ttl time.Duration, scope string, userID int64) (*api.TokenResponse, error) {
	const lengthRandom = 16
	randomBytes := make([]byte, lengthRandom)

	if _, err := rand.Read(randomBytes); err != nil {
		return nil, fmt.Errorf("failed read rand: %w", err)
	}

	plaintext := base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(randomBytes)
	hash := sha256.Sum256([]byte(plaintext))

	token, err := q.CreateToken(ctx, data.CreateTokenParams{Hash: hash[:], UserID: userID, Expiry: time.Now().Add(ttl), Scope: scope})
	if err != nil {
		return nil, fmt.Errorf("failed create token: %w", err)
	}

	tokenPlaintext := &api.TokenResponse{Token: plaintext, Expiry: token.Expiry, Scope: token.Scope}

	return tokenPlaintext, nil
}
