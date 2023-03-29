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

func (s *Handler) NewActivationToken(ctx context.Context, req *api.UserEmailRequest) (api.NewActivationTokenRes, error) {
	activationToken, err := logic.NewActivationToken(ctx, s.Queries, req.Email)
	if err != nil {
		switch {
		case errors.Is(err, logic.ErrEmailNotFound):
			return &api.NewActivationTokenNotFound{Error: err.Error()}, nil
		case errors.Is(err, logic.ErrAlreadyActivated):
			return &api.NewActivationTokenUnprocessableEntity{Error: err.Error()}, nil
		default:
			return &api.NewActivationTokenInternalServerError{Error: serverError}, nil
		}
	}

	// TODO Fix emails
	err = s.Mailer.Send(req.Email, "token_activation.tmpl", map[string]any{
		"activationToken": activationToken.Plaintext,
	})
	if err != nil {
		return &api.NewActivationTokenInternalServerError{Error: serverError}, nil
	}

	return &activationToken, nil
}

func (s *Handler) NewPasswordResetToken(ctx context.Context, req *api.UserEmailRequest) (api.NewPasswordResetTokenRes, error) {
	passwordResetToken, err := logic.NewPasswordResetToken(ctx, s.Queries, req.Email)
	if err != nil {
		switch {
		case errors.Is(err, logic.ErrEmailNotFound):
			return &api.NewPasswordResetTokenNotFound{Error: err.Error()}, nil
		case errors.Is(err, logic.ErrActivationRequired):
			return &api.NewPasswordResetTokenUnprocessableEntity{Error: err.Error()}, nil
		default:
			return &api.NewPasswordResetTokenInternalServerError{Error: serverError}, nil
		}
	}

	err = s.Mailer.Send(req.Email, "token_password_reset.tmpl", map[string]any{
		"passwordResetToken": passwordResetToken.Plaintext,
	})
	if err != nil {
		return &api.NewPasswordResetTokenInternalServerError{Error: serverError}, nil
	}

	return &passwordResetToken, nil
}

func (s *Handler) NewRefreshToken(ctx context.Context, req *api.UserLoginRequest) (api.NewRefreshTokenRes, error) {
	refreshToken, accessToken, err := logic.NewRefreshToken(ctx, s.Queries, req.Email, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, logic.ErrInvalidCredentials):
			return &api.NewRefreshTokenUnauthorized{Error: err.Error()}, nil
		default:
			return &api.NewRefreshTokenInternalServerError{Error: serverError}, nil
		}
	}

	cookie, err := newCookie(cookieRefreshToken, refreshToken.Plaintext, cookieTTL, s.Secret)
	if err != nil {
		return &api.NewRefreshTokenInternalServerError{Error: serverError}, nil
	}

	optString := api.OptString{Value: cookie.String(), Set: true}
	tokenResponseHeaders := api.TokenResponseHeaders{SetCookie: optString, Response: accessToken}

	return &tokenResponseHeaders, nil
}

func (s *Handler) NewAccessToken(ctx context.Context) (api.NewAccessTokenRes, error) {
	value := utils.ContextGetCookieValue(ctx)

	refreshToken, accessToken, err := logic.NewAccessToken(ctx, s.Queries, value)
	if err != nil {
		switch {
		case errors.Is(err, logic.ErrInvalidToken), errors.Is(err, logic.ErrReusedRefreshToken):
			return &api.NewAccessTokenUnauthorized{Error: err.Error()}, nil
		default:
			return &api.NewAccessTokenInternalServerError{Error: serverError}, nil
		}
	}

	cookie, err := newCookie(cookieRefreshToken, refreshToken.Plaintext, cookieTTL, s.Secret)
	if err != nil {
		return &api.NewAccessTokenInternalServerError{Error: serverError}, nil
	}

	optString := api.OptString{Value: cookie.String(), Set: true}
	tokenResponseHeaders := api.TokenResponseHeaders{SetCookie: optString, Response: accessToken}

	return &tokenResponseHeaders, nil
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
