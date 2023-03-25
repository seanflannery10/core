package service

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"fmt"
	"strings"

	"github.com/go-faster/errors"
	"github.com/jackc/pgx/v5"
	"github.com/seanflannery10/core/internal/api"
	"github.com/seanflannery10/core/internal/data"
	"github.com/segmentio/asm/base64"
)

func (s *Handler) NewActivationToken(ctx context.Context, req *api.UserEmailRequest) (r api.NewActivationTokenRes, _ error) {
	activationToken, err := newActivationToken(ctx, s.Queries, req.Email)
	if err != nil {
		return &api.NewActivationTokenInternalServerError{}, nil
	}

	err = s.Mailer.Send(req.Email, "token_activation.tmpl", map[string]any{
		"activationToken": activationToken.Plaintext,
	})
	if err != nil {
		return &api.NewActivationTokenInternalServerError{}, nil
	}

	return &activationToken, nil
}

func newActivationToken(ctx context.Context, q data.Queries, email string) (api.TokenResponse, error) {
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

	activationToken, err := newToken(ctx, q, ttlAcitvationToken, data.ScopeActivation, user.ID)
	if err != nil {
		return api.TokenResponse{}, fmt.Errorf("failed new activation token: %w", err)
	}

	return activationToken, nil
}

func (s *Handler) NewPasswordResetToken(ctx context.Context, req *api.UserEmailRequest) (r api.NewPasswordResetTokenRes, _ error) {
	passwordResetToken, err := newPasswordResetToken(ctx, s.Queries, req.Email)
	if err != nil {
		return &api.NewPasswordResetTokenInternalServerError{}, nil
	}

	err = s.Mailer.Send(req.Email, "token_password_reset.tmpl", map[string]any{
		"passwordResetToken": passwordResetToken.Plaintext,
	})
	if err != nil {
		return &api.NewPasswordResetTokenInternalServerError{}, nil
	}

	return &passwordResetToken, nil
}

func newPasswordResetToken(ctx context.Context, q data.Queries, email string) (api.TokenResponse, error) {
	user, err := q.GetUserFromEmail(ctx, email)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return api.TokenResponse{}, errNotFound
		default:
			return api.TokenResponse{}, fmt.Errorf("failed get user from email (password): %w", err)
		}
	}

	passwordResetToken, err := newToken(ctx, q, ttlPasswordResetToken, data.ScopePasswordReset, user.ID)
	if err != nil {
		return api.TokenResponse{}, err
	}

	return passwordResetToken, nil
}

func (s *Handler) NewRefreshToken(ctx context.Context, req *api.UserLoginRequest) (api.NewRefreshTokenRes, error) {
	refreshToken, accessToken, err := newRefreshToken(ctx, s.Queries, req.Email, req.Password)
	if err != nil {
		return &api.NewRefreshTokenInternalServerError{}, nil
	}

	cookie, err := newCookie(cookieRefreshToken, refreshToken.Plaintext, s.Secret)
	if err != nil {
		return &api.NewRefreshTokenInternalServerError{}, nil
	}

	optString := api.OptString{Value: cookie.Value, Set: true}
	tokenResponseHeaders := api.TokenResponseHeaders{SetCookie: optString, Response: accessToken}

	return &tokenResponseHeaders, nil
}

func newRefreshToken(ctx context.Context, q data.Queries, email, pass string) (refresh, access api.TokenResponse, err error) {
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

func (s *Handler) NewAccessToken(ctx context.Context, params api.NewAccessTokenParams) (r api.NewAccessTokenRes, _ error) {
	encryptedValue, err := base64.URLEncoding.DecodeString(params.CoreRefreshToken)
	if err != nil {
		return &api.NewAccessTokenInternalServerError{}, nil
	}

	block, err := aes.NewCipher(s.Secret)
	if err != nil {
		return &api.NewAccessTokenInternalServerError{}, fmt.Errorf("failed new cipher: %w", err)
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return &api.NewAccessTokenInternalServerError{}, fmt.Errorf("failed new gcm: %w", err)
	}

	nonceSize := aesGCM.NonceSize()

	if len(encryptedValue) < nonceSize {
		return &api.NewAccessTokenInternalServerError{}, nil
	}

	nonce := encryptedValue[:nonceSize]
	ciphertext := encryptedValue[nonceSize:]

	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return &api.NewAccessTokenInternalServerError{}, nil
	}

	_, value, ok := strings.Cut(string(plaintext), ":")
	if !ok {
		return &api.NewAccessTokenInternalServerError{}, nil
	}

	refreshToken, accessToken, err := newAccessToken(ctx, s.Queries, value)
	if err != nil {
		return &api.NewAccessTokenInternalServerError{}, nil
	}

	cookie, err := newCookie(cookieRefreshToken, refreshToken.Plaintext, s.Secret)
	if err != nil {
		return &api.NewAccessTokenInternalServerError{}, nil
	}

	optString := api.OptString{Value: cookie.Value, Set: true}
	tokenResponseHeaders := api.TokenResponseHeaders{SetCookie: optString, Response: accessToken}

	return &tokenResponseHeaders, nil
}

func newAccessToken(ctx context.Context, q data.Queries, plaintext string) (refresh, access api.TokenResponse, err error) {
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
