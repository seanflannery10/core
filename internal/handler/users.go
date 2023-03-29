package handler

import (
	"context"

	"github.com/go-faster/errors"
	"github.com/seanflannery10/core/internal/api"
	"github.com/seanflannery10/core/internal/data"
	"github.com/seanflannery10/core/internal/logic"
	"github.com/seanflannery10/core/internal/shared/mailer"
)

const serverError = "the server encountered a problem and could not process your json"

type Handler struct {
	api.UnimplementedHandler

	Mailer  mailer.Mailer
	Queries *data.Queries
	Secret  []byte
}

func (s *Handler) ActivateUser(ctx context.Context, req *api.TokenRequest) (api.ActivateUserRes, error) {
	user, err := logic.ActivateUser(ctx, s.Queries, req.Token)
	if err != nil {
		switch {
		case errors.Is(err, logic.ErrInvalidToken):
			return &api.ActivateUserUnauthorized{Error: err.Error()}, nil
		default:
			return &api.ActivateUserInternalServerError{Error: serverError}, nil
		}
	}

	return &user, nil
}

func (s *Handler) NewUser(ctx context.Context, req *api.UserRequest) (api.NewUserRes, error) {
	user, refreshToken, err := logic.NewUser(ctx, s.Queries, req.Name, req.Email, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, logic.ErrUserExists):
			return &api.NewUserUnprocessableEntity{Error: err.Error()}, nil
		default:
			return &api.NewUserInternalServerError{Error: serverError}, nil
		}
	}

	err = s.Mailer.Send(user.Email, "token_activation.tmpl", map[string]any{
		"activationToken": refreshToken.Token,
	})
	if err != nil {
		return &api.NewUserInternalServerError{Error: serverError}, nil
	}

	return &user, nil
}

func (s *Handler) UpdateUserPassword(ctx context.Context, req *api.UpdateUserPasswordRequest) (api.UpdateUserPasswordRes, error) {
	acceptanceResponse, err := logic.UpdateUserPassword(ctx, s.Queries, req.Token, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, logic.ErrInvalidToken):
			return &api.UpdateUserPasswordUnauthorized{Error: err.Error()}, nil
		default:
			return &api.UpdateUserPasswordInternalServerError{Error: serverError}, nil
		}
	}

	return &acceptanceResponse, nil
}
