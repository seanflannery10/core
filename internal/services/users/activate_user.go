package users

import (
	"errors"
	"net/http"

	"github.com/go-chi/render"
	"github.com/jackc/pgx/v5"
	"github.com/seanflannery10/core/internal/data"
	"github.com/seanflannery10/core/internal/pkg/errs"
	"github.com/seanflannery10/core/internal/pkg/helpers"
	"github.com/seanflannery10/core/internal/pkg/validator"
	"github.com/seanflannery10/core/internal/services"
)

type activateUserPayload struct {
	TokenPlaintext string `json:"token"`
}

func (p *activateUserPayload) Bind(_ *http.Request) error {
	v := validator.New()

	data.ValidateTokenPlaintext(v, p.TokenPlaintext)

	if v.HasErrors() {
		return validator.NewValidationError(v.Errors)
	}

	return nil
}

// @Summary	activate new inactivate account using a token
// @ID			activate-user
// @Produce	json
// @Success	200	{object}	data.User
// @Router		/users/activate  [patch]
func ActivateUserHandler(env *services.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p := &activateUserPayload{}

		if helpers.CheckAndBind(w, r, p) {
			return
		}

		user, err := env.Queries.GetUserFromTokenHelper(r.Context(), p.TokenPlaintext, data.ScopeActivation)
		if err != nil {
			switch {
			case errors.Is(err, pgx.ErrNoRows):
				_ = render.Render(w, r, errs.ErrInvalidAccessToken())
			default:
				_ = render.Render(w, r, errs.ErrServerError(err))
			}

			return
		}

		user, err = env.Queries.UpdateUser(r.Context(), data.UpdateUserParams{
			UpdateActivated: true,
			Activated:       true,
			ID:              user.ID,
			Version:         user.Version,
		})
		if err != nil {
			_ = render.Render(w, r, errs.ErrServerError(err))
			return
		}

		err = env.Queries.DeleteTokens(r.Context(), data.DeleteTokensParams{
			Scope:  data.ScopeActivation,
			UserID: user.ID,
		})
		if err != nil {
			_ = render.Render(w, r, errs.ErrServerError(err))
		}

		render.Status(r, http.StatusOK)

		helpers.RenderAndCheck(w, r, &user)
	}
}
