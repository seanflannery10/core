package logic

import (
	"context"
	"fmt"

	"github.com/go-faster/errors"
	"github.com/jackc/pgx/v5"
	"github.com/seanflannery10/core/internal/generated/api"
	"github.com/seanflannery10/core/internal/generated/data"
	"golang.org/x/crypto/bcrypt"
)

func ActivateUser(ctx context.Context, q *data.Queries, plaintext string) (*api.UserResponse, error) {
	user, err := getUserFromToken(ctx, q, plaintext, ScopeActivation)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return nil, ErrInvalidToken
		default:
			return nil, fmt.Errorf("failed get user from activation token: %w", err)
		}
	}

	user, err = q.UpdateUser(ctx, data.UpdateUserParams{UpdateActivated: true, Activated: true, ID: user.ID, Version: user.Version})
	if err != nil {
		return nil, fmt.Errorf("failed update user: %w", err)
	}

	if err = q.DeleteTokens(ctx, data.DeleteTokensParams{Scope: ScopeActivation, UserID: user.ID}); err != nil {
		return nil, fmt.Errorf("failed delete tokens: %w", err)
	}

	userResponse := &api.UserResponse{Name: user.Name, Email: user.Email, Version: user.Version}

	return userResponse, nil
}

func NewUser(ctx context.Context, q *data.Queries, name, email, pass string) (*api.UserResponse, *api.TokenResponse, error) {
	ok, err := q.CheckUser(ctx, email)
	if err != nil {
		return nil, nil, fmt.Errorf("failed check user: %w", err)
	}

	if ok {
		return nil, nil, ErrUserExists
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(pass), PasswordCost)
	if err != nil {
		return nil, nil, fmt.Errorf("failed generate password: %w", err)
	}

	passwordHash := hash

	user, err := q.CreateUser(ctx, data.CreateUserParams{Name: name, Email: email, PasswordHash: passwordHash, Activated: false})
	if err != nil {
		return nil, nil, fmt.Errorf("failed create user: %w", err)
	}

	activationToken, err := newToken(ctx, q, ttlActivationToken, ScopeActivation, user.ID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed create new token: %w", err)
	}

	userResponse := &api.UserResponse{Name: user.Name, Email: user.Email, Version: user.Version}

	return userResponse, activationToken, nil
}

func UpdateUserPassword(ctx context.Context, q *data.Queries, token, pass string) (*api.AcceptanceResponse, error) {
	user, err := getUserFromToken(ctx, q, token, ScopePasswordReset)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return nil, ErrInvalidToken
		default:
			return nil, fmt.Errorf("failed get user from password reset token: %w", err)
		}
	}

	user, err = setPassword(user, pass)
	if err != nil {
		return nil, fmt.Errorf("failed set password: %w", err)
	}

	user, err = q.UpdateUser(ctx, data.UpdateUserParams{
		UpdatePasswordHash: true,
		PasswordHash:       user.PasswordHash,
		ID:                 user.ID,
		Version:            user.Version,
	})
	if err != nil {
		return nil, fmt.Errorf("failed update user password: %w", err)
	}

	if err = q.DeleteTokens(ctx, data.DeleteTokensParams{Scope: ScopePasswordReset, UserID: user.ID}); err != nil {
		return nil, fmt.Errorf("failed delete password reset token: %w", err)
	}

	acceptanceResponse := &api.AcceptanceResponse{Message: "password updated"}

	return acceptanceResponse, nil
}
