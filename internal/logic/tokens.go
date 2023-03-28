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
	"github.com/seanflannery10/core/internal/shared/utils"
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
			return api.TokenResponse{}, errNotFound
		default:
			return api.TokenResponse{}, fmt.Errorf("failed get user from email (activation): %w", err)
		}
	}

	if user.Activated {
		return api.TokenResponse{}, errNotActivated
	}

	activationToken, err := utils.NewToken(ctx, q, ttlActivationToken, data.ScopeActivation, user.ID)
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
			return api.TokenResponse{}, errNotFound
		default:
			return api.TokenResponse{}, fmt.Errorf("failed get user from email (password): %w", err)
		}
	}

	passwordResetToken, err := utils.NewToken(ctx, q, ttlPasswordResetToken, data.ScopePasswordReset, user.ID)
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
			return api.TokenResponse{}, api.TokenResponse{}, errNotFound
		default:
			return api.TokenResponse{}, api.TokenResponse{}, fmt.Errorf("failed get user from email (refresh): %w", err)
		}
	}

	match, err := user.ComparePasswords(pass)
	if err != nil {
		return api.TokenResponse{}, api.TokenResponse{}, fmt.Errorf("failed compare passwords: %w", err)
	}

	if !match {
		return api.TokenResponse{}, api.TokenResponse{}, errInvalidCredentials
	}

	refresh, err = utils.NewToken(ctx, q, ttlRefreshToken, data.ScopeRefresh, user.ID)
	if err != nil {
		return api.TokenResponse{}, api.TokenResponse{}, fmt.Errorf("failed create refresh token: %w", err)
	}

	access, err = utils.NewToken(ctx, q, ttlAccessToken, data.ScopeAccess, user.ID)
	if err != nil {
		return api.TokenResponse{}, api.TokenResponse{}, fmt.Errorf("failed create access token: %w", err)
	}

	return refresh, access, nil
}

func NewAccessToken(ctx context.Context, q *data.Queries, plaintext string) (refresh, access api.TokenResponse, err error) {
	user, err := q.GetUserFromTokenHelper(ctx, plaintext, data.ScopeRefresh)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return api.TokenResponse{}, api.TokenResponse{}, errNotFound
		default:
			return api.TokenResponse{}, api.TokenResponse{}, fmt.Errorf("failed get user from refresh token: %w", err)
		}
	}

	tokenHash := sha256.Sum256([]byte(plaintext))

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

		return api.TokenResponse{}, api.TokenResponse{}, errReusedRefreshToken
	}

	err = q.DeactivateToken(ctx, data.DeactivateTokenParams{
		Scope:  data.ScopeRefresh,
		Hash:   tokenHash[:],
		UserID: user.ID,
	})
	if err != nil {
		return api.TokenResponse{}, api.TokenResponse{}, fmt.Errorf("failed deactivate refresh token: %w", err)
	}

	refresh, err = utils.NewToken(ctx, q, ttlRefreshToken, data.ScopeRefresh, user.ID)
	if err != nil {
		return api.TokenResponse{}, api.TokenResponse{}, fmt.Errorf("failed create refresh token: %w", err)
	}

	access, err = utils.NewToken(ctx, q, ttlAccessToken, data.ScopeAccess, user.ID)
	if err != nil {
		return api.TokenResponse{}, api.TokenResponse{}, fmt.Errorf("failed create access token: %w", err)
	}

	return refresh, access, nil
}
