package users

import (
	"net/http"
	"time"

	"github.com/go-chi/render"
	"github.com/seanflannery10/core/internal/data"
	"github.com/seanflannery10/core/pkg/errs"
	"github.com/seanflannery10/core/pkg/helpers"
	"github.com/seanflannery10/core/pkg/httperrors"
	"github.com/seanflannery10/core/pkg/validator"
	"golang.org/x/exp/slog"
)

func createUserHandler(w http.ResponseWriter, r *http.Request) {
	p := &createUserPayload{v: validator.New()}

	err := render.Bind(r, p)
	if err != nil {
		helpers.CheckBindErr(w, r, p.v, err)
		return
	}

	queries := helpers.ContextGetQueries(r)

	user, err := queries.CreateUser(r.Context(), data.CreateUserParams{
		Name:         p.Name,
		Email:        p.Email,
		PasswordHash: p.PasswordHash,
		Activated:    false,
	})
	if err != nil {
		_ = render.Render(w, r, errs.ErrServerError(err))
		return
	}

	token, err := queries.CreateTokenHelper(r.Context(), user.ID, 3*24*time.Hour, data.ScopeActivation)
	if err != nil {
		_ = render.Render(w, r, errs.ErrServerError(err))
		return
	}

	mailer := helpers.ContextGetMailer(r)

	err = mailer.Send(user.Email, "token_activation.tmpl", map[string]any{
		"activationToken": token.Plaintext,
	})
	if err != nil {
		_ = render.Render(w, r, errs.ErrServerError(err))
		return
	}

	render.Status(r, http.StatusCreated)

	err = render.Render(w, r, userResponsePayload{user})
	if err != nil {
		slog.Error("render error", err)
		return
	}
}

func activateUserHandler(w http.ResponseWriter, r *http.Request) {
	p := &activateUserPayload{v: validator.New()}

	err := render.Bind(r, p)
	if err != nil {
		helpers.CheckBindErr(w, r, p.v, err)
		return
	}

	queries := helpers.ContextGetQueries(r)

	user, err := queries.GetUserFromTokenHelper(r.Context(), p.TokenPlaintext, data.ScopeActivation)
	if err != nil {
		_ = render.Render(w, r, errs.ErrServerError(err))
		return
	}

	user, err = queries.UpdateUser(r.Context(), data.UpdateUserParams{
		UpdateActivated: true,
		Activated:       true,
		ID:              user.ID,
		Version:         user.Version,
	})
	if err != nil {
		_ = render.Render(w, r, errs.ErrServerError(err))
		return
	}

	err = queries.DeleteAllTokensForUser(r.Context(), data.DeleteAllTokensForUserParams{
		Scope:  data.ScopeActivation,
		UserID: user.ID,
	})
	if err != nil {
		_ = render.Render(w, r, errs.ErrServerError(err))
	}

	render.Status(r, http.StatusOK)

	err = render.Render(w, r, userResponsePayload{user})
	if err != nil {
		slog.Error("render error", err)
		return
	}
}

func updateUserPasswordHandler(w http.ResponseWriter, r *http.Request) {
	p := &updateUserPasswordPayload{v: validator.New()}

	err := render.Bind(r, p)
	if err != nil {
		helpers.CheckBindErr(w, r, p.v, err)
		return
	}

	queries := helpers.ContextGetQueries(r)

	user, err := queries.GetUserFromTokenHelper(r.Context(), p.TokenPlaintext, data.ScopePasswordReset)
	if err != nil {
		_ = render.Render(w, r, errs.ErrServerError(err))
		return
	}

	err = user.SetPassword(p.Password)
	if err != nil {
		httperrors.ServerError(w, r, err)
		return
	}

	user, err = queries.UpdateUser(r.Context(), data.UpdateUserParams{
		UpdatePasswordHash: true,
		PasswordHash:       user.PasswordHash,
		ID:                 user.ID,
		Version:            user.Version,
	})
	if err != nil {
		_ = render.Render(w, r, errs.ErrServerError(err))
		return
	}

	err = queries.DeleteAllTokensForUser(r.Context(), data.DeleteAllTokensForUserParams{
		Scope:  data.ScopePasswordReset,
		UserID: user.ID,
	})
	if err != nil {
		_ = render.Render(w, r, errs.ErrServerError(err))
	}

	render.Status(r, http.StatusOK)

	err = render.Render(w, r, stringResponsePayload{"your password was successfully reset"})
	if err != nil {
		slog.Error("render error", err)
		return
	}
}