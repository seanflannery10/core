package handler

import (
	"context"
	"fmt"

	"github.com/seanflannery10/core/internal/logic"
	"github.com/seanflannery10/core/internal/oas"
)

func (s *Handler) ActivateUser(ctx context.Context, req *oas.TokenRequest) (oas.ActivateUserRes, error) {
	user, err := logic.ActivateUser(ctx, s.Queries, req.Plaintext)
	if err != nil {
		return &oas.ActivateUserInternalServerError{}, fmt.Errorf("failed handler activate user: %w", err)
	}

	return &user, nil
}

func (s *Handler) NewUser(ctx context.Context, req *oas.UserRequest) (oas.NewUserRes, error) {
	user, refreshToken, err := logic.NewUser(ctx, s.Queries, req.Name, req.Email, req.Password)
	if err != nil {
		return &oas.NewUserInternalServerError{}, fmt.Errorf("failed handler new user: %w", err)
	}

	err = s.Mailer.Send(user.Email, "token_activation.tmpl", map[string]any{
		"activationToken": refreshToken.Plaintext,
	})
	if err != nil {
		return &oas.NewUserInternalServerError{}, nil
	}

	return &user, nil
}

func (s *Handler) UpdateUserPassword(ctx context.Context, req *oas.UpdateUserPasswordRequest) (oas.UpdateUserPasswordRes, error) {
	_, err := logic.UpdateUserPassword(ctx, s.Queries, req.Token, req.Password)
	if err != nil {
		return &oas.UpdateUserPasswordInternalServerError{}, fmt.Errorf("failed handler update user password: %w", err)
	}

	return &oas.AcceptanceResponse{Message: "password updated"}, nil
}
