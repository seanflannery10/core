package users

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/seanflannery10/core/internal/data"
	"github.com/seanflannery10/core/internal/pkg/errs"
	"github.com/seanflannery10/core/internal/pkg/helpers"
	"github.com/seanflannery10/core/internal/pkg/responses"
	"github.com/seanflannery10/core/internal/pkg/validator"
	"github.com/seanflannery10/core/internal/services"
)

type updateUserPasswordPayload struct {
	Password       string `json:"password"`
	TokenPlaintext string `json:"token"`
}

func (p *updateUserPasswordPayload) Bind(_ *http.Request) error {
	v := validator.New()

	data.ValidatePasswordPlaintext(v, p.Password)
	data.ValidateTokenPlaintext(v, p.TokenPlaintext)

	if v.HasErrors() {
		return validator.ErrValidation
	}

	return nil
}

func UpdateUserPasswordHandler(env *services.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p := &updateUserPasswordPayload{}

		if helpers.CheckAndBind(w, r, p) {
			return
		}

		user, err := env.Queries.GetUserFromTokenHelper(r.Context(), p.TokenPlaintext, data.ScopePasswordReset)
		if err != nil {
			_ = render.Render(w, r, errs.ErrServerError(err))
			return
		}

		err = user.SetPassword(p.Password)
		if err != nil {
			_ = render.Render(w, r, errs.ErrServerError(err))
			return
		}

		user, err = env.Queries.UpdateUser(r.Context(), data.UpdateUserParams{
			UpdatePasswordHash: true,
			PasswordHash:       user.PasswordHash,
			ID:                 user.ID,
			Version:            user.Version,
		})
		if err != nil {
			_ = render.Render(w, r, errs.ErrServerError(err))
			return
		}

		err = env.Queries.DeleteTokens(r.Context(), data.DeleteTokensParams{
			Scope:  data.ScopePasswordReset,
			UserID: user.ID,
		})
		if err != nil {
			_ = render.Render(w, r, errs.ErrServerError(err))
		}

		render.Status(r, http.StatusOK)

		helpers.RenderAndCheck(w, r, responses.NewStringResponsePayload("your password was successfully reset"))
	}
}
