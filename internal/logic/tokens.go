package logic

import (
	"context"
	"crypto/sha256"
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

func NewActivationToken(ctx context.Context, q *data.Queries, email string) (*api.TokenResponse, error) {
	user, err := q.GetUserFromEmail(ctx, email)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return nil, ErrEmailNotFound
		default:
			return nil, fmt.Errorf("failed get user from email (activation): %w", err)
		}
	}

	if user.Activated {
		return nil, ErrUserAlreadyActivated
	}

	activationToken, err := newToken(ctx, q, ttlActivationToken, ScopeActivation, user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed new activation token: %w", err)
	}

	return activationToken, nil
}

func NewPasswordResetToken(ctx context.Context, q *data.Queries, email string) (*api.TokenResponse, error) {
	user, err := q.GetUserFromEmail(ctx, email)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return nil, ErrEmailNotFound
		default:
			return nil, fmt.Errorf("failed get user from email (password): %w", err)
		}
	}

	if !user.Activated {
		return nil, ErrActivationRequired
	}

	passwordResetToken, err := newToken(ctx, q, ttlPasswordResetToken, ScopePasswordReset, user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed create password reset token: %w", err)
	}

	return passwordResetToken, nil
}

func NewRefreshToken(ctx context.Context, q *data.Queries, email, pass string) (refresh, access *api.TokenResponse, err error) {
	user, err := q.GetUserFromEmail(ctx, email)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return nil, nil, ErrInvalidCredentials
		default:
			return nil, nil, fmt.Errorf("failed get user from email (refresh): %w", err)
		}
	}

	if err = comparePasswords(user, pass); err != nil {
		return nil, nil, fmt.Errorf("failed compare passwords: %w", err)
	}

	refresh, err = newToken(ctx, q, ttlRefreshToken, ScopeRefresh, user.ID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed create refresh token: %w", err)
	}

	access, err = newToken(ctx, q, ttlAccessToken, ScopeAccess, user.ID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed create access token: %w", err)
	}

	return refresh, access, nil
}

func NewAccessToken(ctx context.Context, q *data.Queries, tokenFromCookie string) (refresh, access *api.TokenResponse, err error) {
	user, err := getUserFromToken(ctx, q, tokenFromCookie, ScopeRefresh)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return nil, nil, ErrInvalidToken
		default:
			return nil, nil, fmt.Errorf("failed get user from refresh token: %w", err)
		}
	}

	tokenHash := sha256.Sum256([]byte(tokenFromCookie))

	badToken, err := q.CheckToken(ctx, data.CheckTokenParams{Hash: tokenHash[:], UserID: user.ID, Scope: ScopeRefresh})
	if err != nil {
		return nil, nil, fmt.Errorf("failed check refresh token: %w", err)
	}

	if badToken {
		if err = q.DeleteTokens(ctx, data.DeleteTokensParams{Scope: ScopeRefresh, UserID: user.ID}); err != nil {
			return nil, nil, fmt.Errorf("failed deactivate refresh token: %w", err)
		}

		return nil, nil, ErrReusedRefreshToken
	}

	if err = q.DeactivateToken(ctx, data.DeactivateTokenParams{Scope: ScopeRefresh, Hash: tokenHash[:], UserID: user.ID}); err != nil {
		return nil, nil, fmt.Errorf("failed deactivate refresh token: %w", err)
	}

	refresh, err = newToken(ctx, q, ttlRefreshToken, ScopeRefresh, user.ID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed create refresh token: %w", err)
	}

	access, err = newToken(ctx, q, ttlAccessToken, ScopeAccess, user.ID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed create access token: %w", err)
	}

	return refresh, access, nil
}
