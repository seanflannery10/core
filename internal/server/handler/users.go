package handler

import (
	"context"

	"github.com/go-faster/errors"
	"github.com/seanflannery10/core/internal/generated/api"
	"github.com/seanflannery10/core/internal/server/logic"
)

func (s *Handler) ActivateUser(ctx context.Context, req *api.TokenRequest) (*api.UserResponse, error) {
	user, err := logic.ActivateUser(ctx, s.Queries, req.Token)
	if err != nil {
		return nil, errors.Wrap(err, "failed activate user")
	}

	return user, nil
}

func (s *Handler) NewUser(ctx context.Context, req *api.UserRequest) (*api.UserResponse, error) {
	user, refreshToken, err := logic.NewUser(ctx, s.Queries, req.Name, req.Email, req.Password)
	if err != nil {
		return nil, errors.Wrap(err, "failed new user")
	}

	if err = s.Mailer.Send(user.Email, "token_activation.tmpl", map[string]any{"activationToken": refreshToken.Token}); err != nil {
		return nil, errors.Wrap(err, "failed send new user activation token email")
	}

	return user, nil
}

func (s *Handler) UpdateUserPassword(ctx context.Context, req *api.UpdateUserPasswordRequest) (*api.AcceptanceResponse, error) {
	acceptanceResponse, err := logic.UpdateUserPassword(ctx, s.Queries, req.Token, req.Password)
	if err != nil {
		return nil, errors.Wrap(err, "failed update user password")
	}

	return acceptanceResponse, nil
}
