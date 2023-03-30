package logic

import (
	"context"
	"fmt"

	"github.com/go-faster/errors"
	"github.com/jackc/pgx/v5"
	"github.com/seanflannery10/core/internal/api"
	"github.com/seanflannery10/core/internal/data"
	"golang.org/x/crypto/bcrypt"
)

func ActivateUser(ctx context.Context, q *data.Queries, plaintext string) (api.UserResponse, error) {
	user, err := q.GetUserFromTokenHelper(ctx, plaintext, data.ScopeActivation)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return api.UserResponse{}, ErrInvalidToken
		default:
			return api.UserResponse{}, fmt.Errorf("failed get user from activation token: %w", err)
		}
	}

	user, err = q.UpdateUser(ctx, data.UpdateUserParams{
		UpdateActivated: true,
		Activated:       true,
		ID:              user.ID,
		Version:         user.Version,
	})
	if err != nil {
		return api.UserResponse{}, fmt.Errorf("failed update user: %w", err)
	}

	err = q.DeleteTokens(ctx, data.DeleteTokensParams{
		Scope:  data.ScopeActivation,
		UserID: user.ID,
	})
	if err != nil {
		return api.UserResponse{}, fmt.Errorf("failed delete tokens: %w", err)
	}

	userResponse := api.UserResponse{
		Name:    user.Name,
		Email:   user.Email,
		Version: user.Version,
	}

	return userResponse, nil
}

func NewUser(ctx context.Context, q *data.Queries, name, email, pass string) (api.UserResponse, api.TokenResponse, error) {
	ok, err := q.CheckUser(ctx, email)
	if err != nil {
		return api.UserResponse{}, api.TokenResponse{}, fmt.Errorf("failed check user: %w", err)
	}

	if ok {
		return api.UserResponse{}, api.TokenResponse{}, ErrUserExists
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(pass), data.PasswordCost)
	if err != nil {
		return api.UserResponse{}, api.TokenResponse{}, fmt.Errorf("failed generate password: %w", err)
	}

	passwordHash := hash

	user, err := q.CreateUser(ctx, data.CreateUserParams{
		Name:         name,
		Email:        email,
		PasswordHash: passwordHash,
		Activated:    false,
	})
	if err != nil {
		return api.UserResponse{}, api.TokenResponse{}, fmt.Errorf("failed create user: %w", err)
	}

	activationToken, err := newToken(ctx, q, ttlActivationToken, data.ScopeActivation, user.ID)
	if err != nil {
		return api.UserResponse{}, api.TokenResponse{}, fmt.Errorf("failed create new token: %w", err)
	}

	userResponse := api.UserResponse{
		Name:    user.Name,
		Email:   user.Email,
		Version: user.Version,
	}

	return userResponse, activationToken, nil
}

func UpdateUserPassword(ctx context.Context, q *data.Queries, token, pass string) (api.AcceptanceResponse, error) {
	user, err := q.GetUserFromTokenHelper(ctx, token, data.ScopePasswordReset)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return api.AcceptanceResponse{}, ErrInvalidToken
		default:
			return api.AcceptanceResponse{}, fmt.Errorf("failed get user from password reset token: %w", err)
		}
	}

	err = user.SetPassword(pass)
	if err != nil {
		return api.AcceptanceResponse{}, fmt.Errorf("failed set password: %w", err)
	}

	user, err = q.UpdateUser(ctx, data.UpdateUserParams{
		UpdatePasswordHash: true,
		PasswordHash:       user.PasswordHash,
		ID:                 user.ID,
		Version:            user.Version,
	})
	if err != nil {
		return api.AcceptanceResponse{}, fmt.Errorf("failed update user password: %w", err)
	}

	err = q.DeleteTokens(ctx, data.DeleteTokensParams{
		Scope:  data.ScopePasswordReset,
		UserID: user.ID,
	})
	if err != nil {
		return api.AcceptanceResponse{}, fmt.Errorf("failed delete password reset token: %w", err)
	}

	return api.AcceptanceResponse{Message: "password updated"}, nil
}
