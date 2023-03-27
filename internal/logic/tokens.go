package logic

import (
	"context"
	"crypto/sha256"
	"fmt"

	"github.com/go-faster/errors"
	"github.com/jackc/pgx/v5"
	"github.com/seanflannery10/core/internal/data"
	"github.com/seanflannery10/core/internal/oas"
	"github.com/seanflannery10/core/internal/shared/utils"
)

func NewActivationToken(ctx context.Context, q data.Queries, email string) (oas.TokenResponse, error) {
	user, err := q.GetUserFromEmail(ctx, email)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return oas.TokenResponse{}, errNotFound
		default:
			return oas.TokenResponse{}, fmt.Errorf("failed get user from email (activation): %w", err)
		}
	}

	if user.Activated {
		return oas.TokenResponse{}, errNotActivated
	}

	activationToken, err := utils.NewToken(ctx, q, ttlActivationToken, data.ScopeActivation, user.ID)
	if err != nil {
		return oas.TokenResponse{}, fmt.Errorf("failed new activation token: %w", err)
	}

	return activationToken, nil
}

func NewPasswordResetToken(ctx context.Context, q data.Queries, email string) (oas.TokenResponse, error) {
	user, err := q.GetUserFromEmail(ctx, email)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return oas.TokenResponse{}, errNotFound
		default:
			return oas.TokenResponse{}, fmt.Errorf("failed get user from email (password): %w", err)
		}
	}

	passwordResetToken, err := utils.NewToken(ctx, q, ttlPasswordResetToken, data.ScopePasswordReset, user.ID)
	if err != nil {
		return oas.TokenResponse{}, fmt.Errorf("failed create password reset token: %w", err)
	}

	return passwordResetToken, nil
}

func NewRefreshToken(ctx context.Context, q data.Queries, email, pass string) (refresh, access oas.TokenResponse, err error) {
	user, err := q.GetUserFromEmail(ctx, email)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return oas.TokenResponse{}, oas.TokenResponse{}, errNotFound
		default:
			return oas.TokenResponse{}, oas.TokenResponse{}, fmt.Errorf("failed get user from email (refresh): %w", err)
		}
	}

	match, err := user.ComparePasswords(pass)
	if err != nil {
		return oas.TokenResponse{}, oas.TokenResponse{}, fmt.Errorf("failed compare passwords: %w", err)
	}

	if !match {
		return oas.TokenResponse{}, oas.TokenResponse{}, errInvalidCredentials
	}

	refresh, err = utils.NewToken(ctx, q, ttlRefreshToken, data.ScopeRefresh, user.ID)
	if err != nil {
		return oas.TokenResponse{}, oas.TokenResponse{}, fmt.Errorf("failed create refresh token: %w", err)
	}

	access, err = utils.NewToken(ctx, q, ttlAccessToken, data.ScopeAccess, user.ID)
	if err != nil {
		return oas.TokenResponse{}, oas.TokenResponse{}, fmt.Errorf("failed create access token: %w", err)
	}

	return refresh, access, nil
}

func NewAccessToken(ctx context.Context, q data.Queries, plaintext string) (refresh, access oas.TokenResponse, err error) {
	user, err := q.GetUserFromTokenHelper(ctx, plaintext, data.ScopeRefresh)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return oas.TokenResponse{}, oas.TokenResponse{}, errNotFound
		default:
			return oas.TokenResponse{}, oas.TokenResponse{}, fmt.Errorf("failed get user from refresh token: %w", err)
		}
	}

	tokenHash := sha256.Sum256([]byte(plaintext))

	badToken, err := q.CheckToken(ctx, data.CheckTokenParams{
		Hash:   tokenHash[:],
		UserID: user.ID,
		Scope:  data.ScopeRefresh,
	})
	if err != nil {
		return oas.TokenResponse{}, oas.TokenResponse{}, fmt.Errorf("failed check refresh token: %w", err)
	}

	if badToken {
		err = q.DeleteTokens(ctx, data.DeleteTokensParams{
			Scope:  data.ScopeRefresh,
			UserID: user.ID,
		})
		if err != nil {
			return oas.TokenResponse{}, oas.TokenResponse{}, fmt.Errorf("failed deactivate refresh token: %w", err)
		}

		return oas.TokenResponse{}, oas.TokenResponse{}, errReusedRefreshToken
	}

	err = q.DeactivateToken(ctx, data.DeactivateTokenParams{
		Scope:  data.ScopeRefresh,
		Hash:   tokenHash[:],
		UserID: user.ID,
	})
	if err != nil {
		return oas.TokenResponse{}, oas.TokenResponse{}, fmt.Errorf("failed deactivate refresh token: %w", err)
	}

	refresh, err = utils.NewToken(ctx, q, ttlRefreshToken, data.ScopeRefresh, user.ID)
	if err != nil {
		return oas.TokenResponse{}, oas.TokenResponse{}, fmt.Errorf("failed create refresh token: %w", err)
	}

	access, err = utils.NewToken(ctx, q, ttlAccessToken, data.ScopeAccess, user.ID)
	if err != nil {
		return oas.TokenResponse{}, oas.TokenResponse{}, fmt.Errorf("failed create access token: %w", err)
	}

	return refresh, access, nil
}
