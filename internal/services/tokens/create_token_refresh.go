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

type createTokenRefreshPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (p *createTokenRefreshPayload) Bind(_ *http.Request) error {
	v := validator.New()

	data.ValidateEmail(v, p.Email)
	data.ValidatePasswordPlaintext(v, p.Password)

	if v.HasErrors() {
		return validator.NewValidationError(v.Errors)
	}

	return nil
}

// @Summary	create refresh token using an email address
// @ID			create-token-refresh
// @Produce	json
// @Success	200	{object}	data.TokenFull
// @Router		/tokens/refresh  [post]
func CreateTokenRefreshHandler(env *services.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p := &createTokenRefreshPayload{}

		if helpers.CheckAndBind(w, r, p) {
			return
		}

		user, err := env.Queries.GetUserFromEmail(r.Context(), p.Email)
		if err != nil {
			switch {
			case errors.Is(err, pgx.ErrNoRows):
				_ = render.Render(w, r, errs.ErrInvalidCredentials())
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
			_ = render.Render(w, r, errs.ErrInvalidCredentials())
			return
		}

		env.User = &user

		w, accessToken, err := createRefreshAndAccessTokens(w, r, env)
		if err != nil {
			_ = render.Render(w, r, errs.ErrServerError(err))
			return
		}

		r = middleware.LogUser(r, &user)

		render.Status(r, http.StatusCreated)

		helpers.RenderAndCheck(w, r, &accessToken)
	}
}
