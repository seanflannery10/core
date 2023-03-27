package handler

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"fmt"
	"strings"

	"github.com/seanflannery10/core/internal/logic"
	"github.com/seanflannery10/core/internal/oas"
	"github.com/seanflannery10/core/internal/shared/utils"
	"github.com/segmentio/asm/base64"
)

func (s *Handler) NewActivationToken(ctx context.Context, req *oas.UserEmailRequest) (oas.NewActivationTokenRes, error) {
	activationToken, err := logic.NewActivationToken(ctx, s.Queries, req.Email)
	if err != nil {
		return &oas.NewActivationTokenInternalServerError{}, nil
	}

	err = s.Mailer.Send(req.Email, "token_activation.tmpl", map[string]any{
		"activationToken": activationToken.Plaintext,
	})
	if err != nil {
		return &oas.NewActivationTokenInternalServerError{}, nil
	}

	return &activationToken, nil
}

func (s *Handler) NewPasswordResetToken(ctx context.Context, req *oas.UserEmailRequest) (oas.NewPasswordResetTokenRes, error) {
	passwordResetToken, err := logic.NewPasswordResetToken(ctx, s.Queries, req.Email)
	if err != nil {
		return &oas.NewPasswordResetTokenInternalServerError{}, nil
	}

	err = s.Mailer.Send(req.Email, "token_password_reset.tmpl", map[string]any{
		"passwordResetToken": passwordResetToken.Plaintext,
	})
	if err != nil {
		return &oas.NewPasswordResetTokenInternalServerError{}, nil
	}

	return &passwordResetToken, nil
}

func (s *Handler) NewRefreshToken(ctx context.Context, req *oas.UserLoginRequest) (oas.NewRefreshTokenRes, error) {
	refreshToken, accessToken, err := logic.NewRefreshToken(ctx, s.Queries, req.Email, req.Password)
	if err != nil {
		return &oas.NewRefreshTokenInternalServerError{}, nil
	}

	cookie, err := utils.NewCookie(cookieRefreshToken, refreshToken.Plaintext, cookieTTL, s.Secret)
	if err != nil {
		return &oas.NewRefreshTokenInternalServerError{}, nil
	}

	tokenResponseHeaders := oas.TokenResponseHeaders{SetCookie: cookie, Response: accessToken}

	return &tokenResponseHeaders, nil
}

func (s *Handler) NewAccessToken(ctx context.Context, params oas.NewAccessTokenParams) (oas.NewAccessTokenRes, error) {
	encryptedValue, err := base64.URLEncoding.DecodeString(params.CoreRefreshToken)
	if err != nil {
		return &oas.NewAccessTokenInternalServerError{}, nil
	}

	block, err := aes.NewCipher(s.Secret)
	if err != nil {
		return &oas.NewAccessTokenInternalServerError{}, fmt.Errorf("failed new cipher: %w", err)
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return &oas.NewAccessTokenInternalServerError{}, fmt.Errorf("failed new gcm: %w", err)
	}

	nonceSize := aesGCM.NonceSize()

	if len(encryptedValue) < nonceSize {
		return &oas.NewAccessTokenInternalServerError{}, nil
	}

	nonce := encryptedValue[:nonceSize]
	ciphertext := encryptedValue[nonceSize:]

	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return &oas.NewAccessTokenInternalServerError{}, nil
	}

	_, value, ok := strings.Cut(string(plaintext), ":")
	if !ok {
		return &oas.NewAccessTokenInternalServerError{}, nil
	}

	refreshToken, accessToken, err := logic.NewAccessToken(ctx, s.Queries, value)
	if err != nil {
		return &oas.NewAccessTokenInternalServerError{}, nil
	}

	cookie, err := utils.NewCookie(cookieRefreshToken, refreshToken.Plaintext, cookieTTL, s.Secret)
	if err != nil {
		return &oas.NewAccessTokenInternalServerError{}, nil
	}

	tokenResponseHeaders := oas.TokenResponseHeaders{SetCookie: cookie, Response: accessToken}

	return &tokenResponseHeaders, nil
}
