package handler

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"fmt"
	"strings"

	"github.com/seanflannery10/core/internal/api"
	"github.com/seanflannery10/core/internal/logic"
	"github.com/seanflannery10/core/internal/shared/utils"
	"github.com/segmentio/asm/base64"
)

const (
	cookieRefreshToken = "core_refreshtoken"
	cookieTTL          = 7 * 24 * 60 * 60
)

func (s *Handler) NewActivationToken(ctx context.Context, req *api.UserEmailRequest) (api.NewActivationTokenRes, error) {
	activationToken, err := logic.NewActivationToken(ctx, s.Queries, req.Email)
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

func (s *Handler) NewPasswordResetToken(ctx context.Context, req *api.UserEmailRequest) (api.NewPasswordResetTokenRes, error) {
	passwordResetToken, err := logic.NewPasswordResetToken(ctx, s.Queries, req.Email)
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

func (s *Handler) NewRefreshToken(ctx context.Context, req *api.UserLoginRequest) (api.NewRefreshTokenRes, error) {
	refreshToken, accessToken, err := logic.NewRefreshToken(ctx, s.Queries, req.Email, req.Password)
	if err != nil {
		return &api.NewRefreshTokenInternalServerError{}, nil
	}

	cookie, err := utils.NewCookie(cookieRefreshToken, refreshToken.Plaintext, cookieTTL, s.Secret)
	if err != nil {
		return &api.NewRefreshTokenInternalServerError{}, nil
	}

	tokenResponseHeaders := api.TokenResponseHeaders{SetCookie: cookie, Response: accessToken}

	return &tokenResponseHeaders, nil
}

func (s *Handler) NewAccessToken(ctx context.Context, params api.NewAccessTokenParams) (api.NewAccessTokenRes, error) {
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

	refreshToken, accessToken, err := logic.NewAccessToken(ctx, s.Queries, value)
	if err != nil {
		return &api.NewAccessTokenInternalServerError{}, nil
	}

	cookie, err := utils.NewCookie(cookieRefreshToken, refreshToken.Plaintext, cookieTTL, s.Secret)
	if err != nil {
		return &api.NewAccessTokenInternalServerError{}, nil
	}

	tokenResponseHeaders := api.TokenResponseHeaders{SetCookie: cookie, Response: accessToken}

	return &tokenResponseHeaders, nil
}
