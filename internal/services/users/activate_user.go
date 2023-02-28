package users

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/seanflannery10/core/internal/data"
	"github.com/seanflannery10/core/pkg/errs"
	"github.com/seanflannery10/core/pkg/helpers"
	"github.com/seanflannery10/core/pkg/validator"
)

type activateUserPayload struct {
	TokenPlaintext string `json:"token"`
}

func (p *activateUserPayload) Bind(_ *http.Request) error {
	v := validator.New()

	data.ValidateTokenPlaintext(v, p.TokenPlaintext)

	if v.HasErrors() {
		return validator.ErrValidation
	}

	return nil
}

func ActivateUserHandler(w http.ResponseWriter, r *http.Request) {
	p := &activateUserPayload{}

	if helpers.CheckAndBind(w, r, p) {
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
