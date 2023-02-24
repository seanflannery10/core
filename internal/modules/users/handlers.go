package users

import (
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/render"
	"github.com/seanflannery10/core/internal/data"
	"github.com/seanflannery10/core/pkg/errs"
	"github.com/seanflannery10/core/pkg/helpers"
	"github.com/seanflannery10/core/pkg/httperrors"
	"golang.org/x/exp/slog"
)

func createUserHandler(w http.ResponseWriter, r *http.Request) {
	input := &Request{}

	err := render.Bind(r, input)
	if err != nil {
		switch {
		case errors.Is(err, errUserExistists):
			_ = render.Render(w, r, errs.ErrInvalidAuthenticationToken)
		default:
			_ = render.Render(w, r, errs.ErrBadRequest(err))
		}

		return
	}

	queries := helpers.ContextGetQueries(r)

	user, err := queries.CreateUser(r.Context(), data.CreateUserParams{
		Name:         input.User.Name,
		Email:        input.User.Email,
		PasswordHash: input.User.PasswordHash,
		Activated:    false,
	})
	if err != nil {
		_ = render.Render(w, r, errs.ErrServerError(err))
		return
	}

	token, err := queries.NewToken(r.Context(), user.ID, 3*24*time.Hour, data.ScopeActivation)
	if err != nil {
		_ = render.Render(w, r, errs.ErrServerError(err))
		return
	}

	slog.Info("debug", "token", token)
	//
	// server.Background(func() {
	//	input := map[string]any{
	//		"activationToken": token.Plaintext,
	//	}
	//
	//	err = mailer.Send(user.Email, "token_activation.tmpl", input)
	//	if err != nil {
	//		slog.Error("email error", err)
	//	}
	// })

	render.Status(r, http.StatusCreated)

	err = render.Render(w, r, NewResponse(user))
	if err != nil {
		httperrors.ServerError(w, r, err)
	}
}
