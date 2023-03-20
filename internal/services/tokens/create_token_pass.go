package tokens

import (
	"errors"
	"net/http"

	"github.com/go-chi/render"
	"github.com/jackc/pgx/v5"
	"github.com/seanflannery10/core/internal/data"
	"github.com/seanflannery10/core/internal/pkg/errs"
	"github.com/seanflannery10/core/internal/pkg/helpers"
	"github.com/seanflannery10/core/internal/pkg/middleware"
	"github.com/seanflannery10/core/internal/pkg/validator"
	"github.com/seanflannery10/core/internal/services"
)

type createTokenPasswordResetPayload struct {
	Email string `json:"email"`
}

func (p *createTokenPasswordResetPayload) Bind(_ *http.Request) error {
	v := validator.New()

	data.ValidateEmail(v, p.Email)

	if v.HasErrors() {
		return validator.NewValidationError(v.Errors)
	}

	return nil
}

func CreateTokenPasswordResetHandler(env *services.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p := &createTokenPasswordResetPayload{}
		v := validator.New()

		if helpers.CheckAndBind(w, r, p) {
			return
		}

		user, err := env.Queries.GetUserFromEmail(r.Context(), p.Email)
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

		token, err := env.Queries.CreateTokenHelper(r.Context(), user.ID, ttlPasswordResetToken, data.ScopePasswordReset)
		if err != nil {
			_ = render.Render(w, r, errs.ErrServerError(err))
			return
		}

		err = env.Mailer.Send(user.Email, "token_password_reset.tmpl", map[string]any{
			"passwordResetToken": token.Plaintext,
		})
		if err != nil {
			_ = render.Render(w, r, errs.ErrServerError(err))
			return
		}

		r = middleware.LogUser(r, &user)

		render.Status(r, http.StatusCreated)

		helpers.RenderAndCheck(w, r, &token)
	}
}
