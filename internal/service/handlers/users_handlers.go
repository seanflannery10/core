package handlers

import (
	"context"
	"fmt"

	"github.com/seanflannery10/core/internal/api"
)

func (s *Service) ActivateUser(ctx context.Context, req *api.TokenRequest) (api.ActivateUserRes, error) {
	user, err := activateUser(ctx, s.queries, req.Plaintext)
	if err != nil {
		return &api.ActivateUserInternalServerError{}, fmt.Errorf("failed handler activate user: %w", err)
	}

	return &user, nil
}

func (s *Service) NewUser(ctx context.Context, req *api.UserRequest) (api.NewUserRes, error) {
	user, _, err := newUser(ctx, s.queries, req.Name, req.Email, req.Password)
	if err != nil {
		return &api.NewUserInternalServerError{}, fmt.Errorf("failed handler new user: %w", err)
	}

	// w, err = newCookie(w, cookieRefreshToken, refreshToken.Plaintext, t.secret)
	// if err != nil {
	//	return &api.NewRefreshTokenInternalServerError{}, nil
	// }

	return &user, nil
}

func (s *Service) UpdateUserPassword(ctx context.Context, req *api.UpdateUserPasswordRequest) (api.UpdateUserPasswordRes, error) {
	_, err := updateUserPassword(ctx, s.queries, req.Token, req.Password)
	if err != nil {
		return &api.UpdateUserPasswordInternalServerError{}, fmt.Errorf("failed handler update user password: %w", err)
	}

	message := api.AcceptanceResponse{Message: "password updated"}

	return &message, nil
}
