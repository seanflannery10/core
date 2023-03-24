package handlers

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/seanflannery10/core/internal/api"
)

func (s *Service) NewActivationToken(ctx context.Context, req *api.UserEmailRequest) (r api.NewActivationTokenRes, _ error) {
	activationToken, err := newActivationToken(ctx, s.queries, req.Email)
	if err != nil {
		return &api.NewActivationTokenInternalServerError{}, nil
	}

	err = s.mailer.Send(req.Email, "token_activation.tmpl", map[string]any{
		"activationToken": activationToken.Plaintext,
	})
	if err != nil {
		return &api.NewActivationTokenInternalServerError{}, nil
	}

	return &activationToken, nil
}

func (s *Service) NewPasswordResetToken(ctx context.Context, req *api.UserEmailRequest) (r api.NewPasswordResetTokenRes, _ error) {
	passwordResetToken, err := newPasswordResetToken(ctx, s.queries, req.Email)
	if err != nil {
		return &api.NewPasswordResetTokenInternalServerError{}, nil
	}

	err = s.mailer.Send(req.Email, "token_password_reset.tmpl", map[string]any{
		"passwordResetToken": passwordResetToken.Plaintext,
	})
	if err != nil {
		return &api.NewPasswordResetTokenInternalServerError{}, nil
	}

	return &passwordResetToken, nil
}

func (s *Service) NewRefreshToken(ctx context.Context, req *api.UserLoginRequest) (r api.NewRefreshTokenRes, _ error) {
	_, accessToken, err := newRefreshToken(ctx, s.queries, req.Email, req.Password)
	if err != nil {
		return &api.NewRefreshTokenInternalServerError{}, nil
	}

	// w, err = newCookie(w, cookieRefreshToken, refreshToken.Plaintext, t.secret)
	// if err != nil {
	//	return &api.NewRefreshTokenInternalServerError{}, nil
	//}

	return &accessToken, nil
}

func (s *Service) NewAccessToken(ctx context.Context, params api.NewAccessTokenParams) (r api.NewAccessTokenRes, _ error) {
	encryptedValue, err := base64.URLEncoding.DecodeString(params.CoreRefreshToken)
	if err != nil {
		return &api.NewAccessTokenInternalServerError{}, nil
	}

	block, err := aes.NewCipher(s.secret)
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

	_, accessToken, err := newAccessToken(ctx, s.queries, value)
	if err != nil {
		return &api.NewAccessTokenInternalServerError{}, nil
	}

	// w, err = newCookie(w, cookieRefreshToken, refreshToken.Plaintext, t.secret)
	// if err != nil {
	//	return &api.NewAccessTokenInternalServerError{}, nil
	// }

	return &accessToken, nil
}
