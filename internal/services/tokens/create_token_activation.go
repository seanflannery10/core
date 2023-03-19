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

type createTokenActivationPayload struct {
	Email string `json:"email"`
}

func (p *createTokenActivationPayload) Bind(_ *http.Request) error {
	v := validator.New()

	data.ValidateEmail(v, p.Email)

	if v.HasErrors() {
		return validator.NewValidationError(v.Errors)
	}

	return nil
}

// @Summary	create activation token using an email address
// @ID			create-token-activation
// @Produce	json
// @Success	200	{object}	data.TokenFull
// @Router		/tokens/activation  [post]
func CreateTokenActivationHandler(env *services.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p := &createTokenActivationPayload{}
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

		if user.Activated {
			v := validator.New()
			v.AddError("email", "user has already been activated")
			_ = render.Render(w, r, errs.ErrFailedValidation(v.Errors))

			return
		}

		token, err := env.Queries.CreateTokenHelper(r.Context(), user.ID, ttlAcitvationToken, data.ScopeActivation)
		if err != nil {
			_ = render.Render(w, r, errs.ErrServerError(err))
			return
		}

		err = env.Mailer.Send(user.Email, "token_activation.tmpl", map[string]any{
			"activationToken": token.Plaintext,
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
