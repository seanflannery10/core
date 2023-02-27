package tokens

import (
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/render"
	"github.com/jackc/pgx/v5"
	"github.com/seanflannery10/core/internal/data"
	"github.com/seanflannery10/core/pkg/errs"
	"github.com/seanflannery10/core/pkg/helpers"
	"github.com/seanflannery10/core/pkg/httperrors"
	"github.com/seanflannery10/core/pkg/validator"
	"golang.org/x/exp/slog"
)

func createAuthTokenHandler(w http.ResponseWriter, r *http.Request) {
	p := &createAuthTokenPayload{v: validator.New()}

	err := render.Bind(r, p)
	if err != nil {
		helpers.CheckBindErr(w, r, p.v, err)
		return
	}

	q := helpers.ContextGetQueries(r)

	user, err := q.GetUserFromEmail(r.Context(), p.Email)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			httperrors.InvalidCredentials(w, r)
		default:
			httperrors.ServerError(w, r, err)
		}

		return
	}

	match, err := user.ComparePasswords(p.Password)
	if err != nil {
		httperrors.ServerError(w, r, err)
		return
	}

	if !match {
		httperrors.InvalidCredentials(w, r)
		return
	}

	token, err := q.CreateTokenHelper(r.Context(), user.ID, 3*24*time.Hour, data.ScopeAuthentication)
	if err != nil {
		httperrors.ServerError(w, r, err)
		return
	}

	render.Status(r, http.StatusCreated)

	err = render.Render(w, r, tokenResponsePayload{token})
	if err != nil {
		slog.Error("render error", err)
	}
}

func createPasswordResetTokenHandler(w http.ResponseWriter, r *http.Request) {
	p := &createPasswordResetTokenPayload{v: validator.New()}

	err := render.Bind(r, p)
	if err != nil {
		helpers.CheckBindErr(w, r, p.v, err)
		return
	}

	q := helpers.ContextGetQueries(r)

	user, err := q.GetUserFromEmail(r.Context(), p.Email)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			p.v.AddError("email", "no matching email address found")
			httperrors.FailedValidation(w, r, p.v)
		default:
			httperrors.ServerError(w, r, err)
		}

		return
	}

	if !user.Activated {
		p.v.AddError("email", "user account must be activated")
		httperrors.FailedValidation(w, r, p.v)

		return
	}

	token, err := q.CreateTokenHelper(r.Context(), user.ID, 45*time.Minute, data.ScopePasswordReset)
	if err != nil {
		httperrors.ServerError(w, r, err)
		return
	}

	mailer := helpers.ContextGetMailer(r)

	err = mailer.Send(user.Email, "token_password_reset.tmpl", map[string]any{
		"passwordResetToken": token.Plaintext,
	})
	if err != nil {
		_ = render.Render(w, r, errs.ErrServerError(err))
		return
	}

	render.Status(r, http.StatusCreated)

	err = render.Render(w, r, tokenResponsePayload{token})
	if err != nil {
		slog.Error("render error", err)
	}
}

type createActivationTokenPayload struct {
	Email string `json:"email"`
	v     *validator.Validator
}

func (p *createActivationTokenPayload) Bind(r *http.Request) error {
	data.ValidateEmail(p.v, p.Email)

	if p.v.HasErrors() {
		return validator.ErrValidation
	}

	return nil
}

func createActivationTokenHandler(w http.ResponseWriter, r *http.Request) {
	p := &createActivationTokenPayload{v: validator.New()}

	err := render.Bind(r, p)
	if err != nil {
		helpers.CheckBindErr(w, r, p.v, err)
		return
	}

	q := helpers.ContextGetQueries(r)

	user, err := q.GetUserFromEmail(r.Context(), p.Email)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			p.v.AddError("email", "no matching email address found")
			httperrors.FailedValidation(w, r, p.v)
		default:
			httperrors.ServerError(w, r, err)
		}

		return
	}

	if user.Activated {
		p.v.AddError("email", "user has already been activated")
		httperrors.FailedValidation(w, r, p.v)

		return
	}

	token, err := q.CreateTokenHelper(r.Context(), user.ID, 3*24*time.Hour, data.ScopeActivation)
	if err != nil {
		httperrors.ServerError(w, r, err)
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

	err = render.Render(w, r, tokenResponsePayload{})
	if err != nil {
		slog.Error("render error", err)
	}
}
