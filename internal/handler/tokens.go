package handler

import (
	"context"

	"github.com/go-faster/errors"
	"github.com/seanflannery10/core/internal/api"
	"github.com/seanflannery10/core/internal/logic"
	"github.com/seanflannery10/core/internal/shared/utils"
)

func (s *Handler) NewActivationToken(ctx context.Context, req *api.UserEmailRequest) (*api.TokenResponse, error) {
	activationToken, err := logic.NewActivationToken(ctx, s.Queries, req.Email)
	if err != nil {
		return nil, errors.Wrap(err, "failed new activation token")
	}

	if err = s.Mailer.Send(req.Email, "token_activation.tmpl", map[string]any{"activationToken": activationToken.Token}); err != nil {
		return nil, errors.Wrap(err, "failed send activation token email")
	}

	return activationToken, nil
}

func (s *Handler) NewPasswordResetToken(ctx context.Context, req *api.UserEmailRequest) (*api.TokenResponse, error) {
	passwordResetToken, err := logic.NewPasswordResetToken(ctx, s.Queries, req.Email)
	if err != nil {
		return nil, errors.Wrap(err, "failed new password reset token")
	}

	if err = s.Mailer.Send(req.Email, "token_password_reset.tmpl", map[string]any{"passwordResetToken": passwordResetToken.Token}); err != nil {
		return nil, errors.Wrap(err, "failed send password reset token email")
	}

	return passwordResetToken, nil
}

func (s *Handler) NewRefreshToken(ctx context.Context, req *api.UserLoginRequest) (*api.TokenResponseHeaders, error) {
	refreshToken, accessToken, err := logic.NewRefreshToken(ctx, s.Queries, req.Email, req.Password)
	if err != nil {
		return nil, errors.Wrap(err, "failed new refresh token")
	}

	cookie, err := newCookie(cookieRefreshToken, refreshToken.Token, cookieTTL, s.Secret)
	if err != nil {
		return nil, errors.Wrap(err, "failed new refresh token cookie")
	}

	optString := api.OptString{Value: cookie.String(), Set: true}
	tokenResponseHeaders := &api.TokenResponseHeaders{SetCookie: optString, Response: *accessToken}

	return tokenResponseHeaders, nil
}

func (s *Handler) NewAccessToken(ctx context.Context) (*api.TokenResponseHeaders, error) {
	value := utils.ContextGetCookieValue(ctx)

	refreshToken, accessToken, err := logic.NewAccessToken(ctx, s.Queries, value)
	if err != nil {
		return nil, errors.Wrap(err, "failed new access token")
	}

	cookie, err := newCookie(cookieRefreshToken, refreshToken.Token, cookieTTL, s.Secret)
	if err != nil {
		return nil, errors.Wrap(err, "failed new access token cookie")
	}

	optString := api.OptString{Value: cookie.String(), Set: true}
	tokenResponseHeaders := &api.TokenResponseHeaders{SetCookie: optString, Response: *accessToken}

	return tokenResponseHeaders, err
}
