package logic

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
	"fmt"
	"time"

	"github.com/go-faster/errors"
	"github.com/jackc/pgx/v5"
	"github.com/seanflannery10/core/internal/api"
	"github.com/seanflannery10/core/internal/data"
)

const (
	ttlAccessToken        = time.Hour
	ttlActivationToken    = 3 * 24 * time.Hour
	ttlPasswordResetToken = 45 * time.Minute
	ttlRefreshToken       = 7 * 24 * time.Hour
)

func NewActivationToken(ctx context.Context, q *data.Queries, email string) (api.TokenResponse, error) {
	user, err := q.GetUserFromEmail(ctx, email)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return api.TokenResponse{}, ErrEmailNotFound
		default:
			return api.TokenResponse{}, fmt.Errorf("failed get user from email (activation): %w", err)
		}
	}

	if user.Activated {
		return api.TokenResponse{}, ErrUserAlreadyActivated
	}

	activationToken, err := newToken(ctx, q, ttlActivationToken, data.ScopeActivation, user.ID)
	if err != nil {
		return api.TokenResponse{}, fmt.Errorf("failed new activation token: %w", err)
	}

	return activationToken, nil
}

func NewPasswordResetToken(ctx context.Context, q *data.Queries, email string) (api.TokenResponse, error) {
	user, err := q.GetUserFromEmail(ctx, email)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return api.TokenResponse{}, ErrEmailNotFound
		default:
			return api.TokenResponse{}, fmt.Errorf("failed get user from email (password): %w", err)
		}
	}

	if !user.Activated {
		return api.TokenResponse{}, ErrActivationRequired
	}

	passwordResetToken, err := newToken(ctx, q, ttlPasswordResetToken, data.ScopePasswordReset, user.ID)
	if err != nil {
		return api.TokenResponse{}, fmt.Errorf("failed create password reset token: %w", err)
	}

	return passwordResetToken, nil
}

func NewRefreshToken(ctx context.Context, q *data.Queries, email, pass string) (refresh, access api.TokenResponse, err error) {
	user, err := q.GetUserFromEmail(ctx, email)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return api.TokenResponse{}, api.TokenResponse{}, ErrInvalidCredentials
		default:
			return api.TokenResponse{}, api.TokenResponse{}, fmt.Errorf("failed get user from email (refresh): %w", err)
		}
	}

	match, err := user.ComparePasswords(pass)
	if err != nil {
		return api.TokenResponse{}, api.TokenResponse{}, fmt.Errorf("failed compare passwords: %w", err)
	}

	if !match {
		return api.TokenResponse{}, api.TokenResponse{}, ErrInvalidCredentials
	}

	refresh, err = newToken(ctx, q, ttlRefreshToken, data.ScopeRefresh, user.ID)
	if err != nil {
		return api.TokenResponse{}, api.TokenResponse{}, fmt.Errorf("failed create refresh token: %w", err)
	}

	access, err = newToken(ctx, q, ttlAccessToken, data.ScopeAccess, user.ID)
	if err != nil {
		return api.TokenResponse{}, api.TokenResponse{}, fmt.Errorf("failed create access token: %w", err)
	}

	return refresh, access, nil
}

func NewAccessToken(ctx context.Context, q *data.Queries, tokenFromCookie string) (refresh, access api.TokenResponse, err error) {
	user, err := q.GetUserFromTokenHelper(ctx, tokenFromCookie, data.ScopeRefresh)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return api.TokenResponse{}, api.TokenResponse{}, ErrInvalidToken
		default:
			return api.TokenResponse{}, api.TokenResponse{}, fmt.Errorf("failed get user from refresh token: %w", err)
		}
	}

	tokenHash := sha256.Sum256([]byte(tokenFromCookie))

	badToken, err := q.CheckToken(ctx, data.CheckTokenParams{
		Hash:   tokenHash[:],
		UserID: user.ID,
		Scope:  data.ScopeRefresh,
	})
	if err != nil {
		return api.TokenResponse{}, api.TokenResponse{}, fmt.Errorf("failed check refresh token: %w", err)
	}

	if badToken {
		err = q.DeleteTokens(ctx, data.DeleteTokensParams{
			Scope:  data.ScopeRefresh,
			UserID: user.ID,
		})
		if err != nil {
			return api.TokenResponse{}, api.TokenResponse{}, fmt.Errorf("failed deactivate refresh token: %w", err)
		}

		return api.TokenResponse{}, api.TokenResponse{}, ErrReusedRefreshToken
	}

	err = q.DeactivateToken(ctx, data.DeactivateTokenParams{
		Scope:  data.ScopeRefresh,
		Hash:   tokenHash[:],
		UserID: user.ID,
	})
	if err != nil {
		return api.TokenResponse{}, api.TokenResponse{}, fmt.Errorf("failed deactivate refresh token: %w", err)
	}

	refresh, err = newToken(ctx, q, ttlRefreshToken, data.ScopeRefresh, user.ID)
	if err != nil {
		return api.TokenResponse{}, api.TokenResponse{}, fmt.Errorf("failed create refresh token: %w", err)
	}

	access, err = newToken(ctx, q, ttlAccessToken, data.ScopeAccess, user.ID)
	if err != nil {
		return api.TokenResponse{}, api.TokenResponse{}, fmt.Errorf("failed create access token: %w", err)
	}

	return refresh, access, nil
}

func newToken(ctx context.Context, q *data.Queries, ttl time.Duration, scope string, userID int64) (api.TokenResponse, error) {
	const lengthRandom = 16
	randomBytes := make([]byte, lengthRandom)

	_, err := rand.Read(randomBytes)
	if err != nil {
		return api.TokenResponse{}, fmt.Errorf("failed read rand: %w", err)
	}

	plaintext := base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(randomBytes)
	hash := sha256.Sum256([]byte(plaintext))

	token, err := q.CreateToken(ctx, data.CreateTokenParams{
		Hash:   hash[:],
		UserID: userID,
		Expiry: time.Now().Add(ttl),
		Scope:  scope,
	})
	if err != nil {
		return api.TokenResponse{}, fmt.Errorf("failed create token: %w", err)
	}

	tokenPlaintext := api.TokenResponse{
		Token:  plaintext,
		Expiry: api.OptDateTime{Value: token.Expiry, Set: true},
		Scope:  api.OptString{Value: token.Scope, Set: true},
	}

	return tokenPlaintext, nil
}
