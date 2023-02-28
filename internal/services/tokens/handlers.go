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
	"github.com/seanflannery10/core/pkg/validator"
)

func createAuthTokenHandler(w http.ResponseWriter, r *http.Request) {
	p := &createAuthTokenPayload{}

	if helpers.CheckAndBind(w, r, p) {
		return
	}

	q := helpers.ContextGetQueries(r)

	user, err := q.GetUserFromEmail(r.Context(), p.Email)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			_ = render.Render(w, r, errs.ErrInvalidCredentials)
		default:
			_ = render.Render(w, r, errs.ErrServerError(err))
		}

		return
	}

	match, err := user.ComparePasswords(p.Password)
	if err != nil {
		_ = render.Render(w, r, errs.ErrServerError(err))
		return
	}

	if !match {
		_ = render.Render(w, r, errs.ErrInvalidCredentials)
		return
	}

	token, err := q.CreateTokenHelper(r.Context(), user.ID, 3*24*time.Hour, data.ScopeAuthentication)
	if err != nil {
		_ = render.Render(w, r, errs.ErrServerError(err))
		return
	}

	render.Status(r, http.StatusCreated)

	helpers.RenderAndCheck(w, r, &token)
}

func createPasswordResetTokenHandler(w http.ResponseWriter, r *http.Request) {
	p := &createPasswordResetTokenPayload{}
	v := validator.New()

	if helpers.CheckAndBind(w, r, p) {
		return
	}

	q := helpers.ContextGetQueries(r)

	user, err := q.GetUserFromEmail(r.Context(), p.Email)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			v.AddError("email", "no matching email address found")
			_ = render.Render(w, r, errs.ErrFailedValidation(v.Errors))
		default:
			_ = render.Render(w, r, errs.ErrServerError(err))
		}

		return
	}

	if !user.Activated {
		v.AddError("email", "user account must be activated")
		_ = render.Render(w, r, errs.ErrFailedValidation(v.Errors))

		return
	}

	token, err := q.CreateTokenHelper(r.Context(), user.ID, 45*time.Minute, data.ScopePasswordReset)
	if err != nil {
		_ = render.Render(w, r, errs.ErrServerError(err))
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

	helpers.RenderAndCheck(w, r, &token)
}

func createActivationTokenHandler(w http.ResponseWriter, r *http.Request) {
	p := &createActivationTokenPayload{}
	v := validator.New()

	if helpers.CheckAndBind(w, r, p) {
		return
	}

	q := helpers.ContextGetQueries(r)

	user, err := q.GetUserFromEmail(r.Context(), p.Email)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			v.AddError("email", "no matching email address found")
			_ = render.Render(w, r, errs.ErrFailedValidation(v.Errors))
		default:
			_ = render.Render(w, r, errs.ErrServerError(err))
		}

		return
	}

	if user.Activated {
		v := validator.New()
		v.AddError("email", "user has already been activated")
		_ = render.Render(w, r, errs.ErrFailedValidation(v.Errors))

		return
	}

	token, err := q.CreateTokenHelper(r.Context(), user.ID, 3*24*time.Hour, data.ScopeActivation)
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

	helpers.RenderAndCheck(w, r, &token)
}
