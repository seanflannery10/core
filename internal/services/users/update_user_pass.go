package users

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/seanflannery10/core/internal/data"
	"github.com/seanflannery10/core/pkg/errs"
	"github.com/seanflannery10/core/pkg/helpers"
	"github.com/seanflannery10/core/pkg/responses"
	"github.com/seanflannery10/core/pkg/validator"
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

func UpdateUserPasswordHandler(w http.ResponseWriter, r *http.Request) {
	p := &updateUserPasswordPayload{}

	if helpers.CheckAndBind(w, r, p) {
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
