package users

import (
	"net/http"
	"time"

	"github.com/go-chi/render"
	"github.com/seanflannery10/core/internal/data"
	"github.com/seanflannery10/core/pkg/errs"
	"github.com/seanflannery10/core/pkg/helpers"
	"github.com/seanflannery10/core/pkg/responses"
	"github.com/seanflannery10/core/pkg/validator"
)

func createUserHandler(w http.ResponseWriter, r *http.Request) {
	p := &createUserPayload{v: validator.New()}

	if helpers.CheckAndBind(w, r, p, p.v) {
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

	helpers.RenderAndCheck(w, r, &user)
}

func activateUserHandler(w http.ResponseWriter, r *http.Request) {
	p := &activateUserPayload{v: validator.New()}

	if helpers.CheckAndBind(w, r, p, p.v) {
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

	helpers.RenderAndCheck(w, r, &user)
}

func updateUserPasswordHandler(w http.ResponseWriter, r *http.Request) {
	p := &updateUserPasswordPayload{v: validator.New()}

	if helpers.CheckAndBind(w, r, p, p.v) {
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
		_ = render.Render(w, r, errs.ErrServerError(err))
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

	helpers.RenderAndCheck(w, r, responses.NewStringResponsePayload("your password was successfully reset"))
}
