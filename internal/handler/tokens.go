package handler

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
	"net/http"

	"github.com/go-faster/errors"
	"github.com/seanflannery10/core/internal/api"
	"github.com/seanflannery10/core/internal/logic"
	"github.com/seanflannery10/core/internal/shared/utils"
	"github.com/segmentio/asm/base64"
)

var errValueTooLong = errors.New("cookie value too long")

const (
	cookieMaxSize      = 4096
	cookieRefreshToken = "core_refresh_token"
	cookieTTL          = 7 * 24 * 60 * 60
)

func (s *Handler) NewActivationToken(ctx context.Context, req *api.UserEmailRequest) (*api.TokenResponse, error) {
	activationToken, err := logic.NewActivationToken(ctx, s.Queries, req.Email)
	if err != nil {
		return &api.TokenResponse{}, errors.Wrap(err, "failed new activation token")
	}

	err = s.Mailer.Send(req.Email, "token_activation.tmpl", map[string]any{
		"activationToken": activationToken.Token,
	})
	if err != nil {
		return &api.TokenResponse{}, errors.Wrap(err, "failed send activation token email")
	}

	return &activationToken, nil
}

func (s *Handler) NewPasswordResetToken(ctx context.Context, req *api.UserEmailRequest) (*api.TokenResponse, error) {
	passwordResetToken, err := logic.NewPasswordResetToken(ctx, s.Queries, req.Email)
	if err != nil {
		return &api.TokenResponse{}, errors.Wrap(err, "failed new password reset token")
	}

	err = s.Mailer.Send(req.Email, "token_password_reset.tmpl", map[string]any{
		"passwordResetToken": passwordResetToken.Token,
	})
	if err != nil {
		return &api.TokenResponse{}, errors.Wrap(err, "failed send password reset token email")
	}

	return &passwordResetToken, nil
}

func (s *Handler) NewRefreshToken(ctx context.Context, req *api.UserLoginRequest) (*api.TokenResponseHeaders, error) {
	refreshToken, accessToken, err := logic.NewRefreshToken(ctx, s.Queries, req.Email, req.Password)
	if err != nil {
		return &api.TokenResponseHeaders{}, errors.Wrap(err, "failed new refresh token")
	}

	cookie, err := newCookie(cookieRefreshToken, refreshToken.Token, cookieTTL, s.Secret)
	if err != nil {
		return &api.TokenResponseHeaders{}, errors.Wrap(err, "failed new refresh token cookie")
	}

	optString := api.OptString{Value: cookie.String(), Set: true}
	tokenResponseHeaders := api.TokenResponseHeaders{SetCookie: optString, Response: accessToken}

	return &tokenResponseHeaders, nil
}

func (s *Handler) NewAccessToken(ctx context.Context) (*api.TokenResponseHeaders, error) {
	value := utils.ContextGetCookieValue(ctx)

	refreshToken, accessToken, err := logic.NewAccessToken(ctx, s.Queries, value)
	if err != nil {
		return &api.TokenResponseHeaders{}, errors.Wrap(err, "failed new access token")
	}

	cookie, err := newCookie(cookieRefreshToken, refreshToken.Token, cookieTTL, s.Secret)
	if err != nil {
		return &api.TokenResponseHeaders{}, errors.Wrap(err, "failed new access token cookie")
	}

	optString := api.OptString{Value: cookie.String(), Set: true}
	tokenResponseHeaders := api.TokenResponseHeaders{SetCookie: optString, Response: accessToken}

	return &tokenResponseHeaders, err
}

func newCookie(name, value string, ttl int, secret []byte) (http.Cookie, error) {
	cookie := http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		MaxAge:   ttl,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}

	block, err := aes.NewCipher(secret)
	if err != nil {
		return http.Cookie{}, fmt.Errorf("failed new cipher: %w", err)
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return http.Cookie{}, fmt.Errorf("failed new gcm: %w", err)
	}

	nonce := make([]byte, aesGCM.NonceSize())

	_, err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return http.Cookie{}, fmt.Errorf("failed read full: %w", err)
	}

	plaintext := fmt.Sprintf("%s:%s", cookie.Name, cookie.Value)

	encryptedValue := aesGCM.Seal(nonce, nonce, []byte(plaintext), nil)

	cookie.Value = base64.URLEncoding.EncodeToString(encryptedValue)

	if len(cookie.String()) > cookieMaxSize {
		return http.Cookie{}, fmt.Errorf("failed length check: %w", errValueTooLong)
	}

	return cookie, nil
}
