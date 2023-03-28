package handler

import (
	"context"
	"fmt"

	"github.com/seanflannery10/core/internal/api"
	"github.com/seanflannery10/core/internal/data"
	"github.com/seanflannery10/core/internal/logic"
	"github.com/seanflannery10/core/internal/shared/mailer"
)

type Handler struct {
	api.UnimplementedHandler

	Mailer  mailer.Mailer
	Queries *data.Queries
	Secret  []byte
}

func (s *Handler) ActivateUser(ctx context.Context, req *api.TokenRequest) (api.ActivateUserRes, error) {
	user, err := logic.ActivateUser(ctx, s.Queries, req.Plaintext)
	if err != nil {
		return &api.ActivateUserInternalServerError{}, fmt.Errorf("failed handler activate user: %w", err)
	}

	return &user, nil
}

func (s *Handler) NewUser(ctx context.Context, req *api.UserRequest) (api.NewUserRes, error) {
	user, refreshToken, err := logic.NewUser(ctx, s.Queries, req.Name, req.Email, req.Password)
	if err != nil {
		return &api.NewUserInternalServerError{Message: "tes", Error: err.Error()}, fmt.Errorf("failed handler new user: %w", err)
	}

	err = s.Mailer.Send(user.Email, "token_activation.tmpl", map[string]any{
		"activationToken": refreshToken.Plaintext,
	})
	if err != nil {
		return &api.NewUserInternalServerError{}, nil
	}

	return &user, nil
}

func (s *Handler) UpdateUserPassword(ctx context.Context, req *api.UpdateUserPasswordRequest) (api.UpdateUserPasswordRes, error) {
	_, err := logic.UpdateUserPassword(ctx, s.Queries, req.Token, req.Password)
	if err != nil {
		return &api.UpdateUserPasswordInternalServerError{}, fmt.Errorf("failed handler update user password: %w", err)
	}

	return &api.AcceptanceResponse{Message: "password updated"}, nil
}
