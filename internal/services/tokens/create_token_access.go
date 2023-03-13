package tokens

import (
	"net/http"
	"time"

	"github.com/go-chi/render"
	"github.com/seanflannery10/core/internal/data"
	"github.com/seanflannery10/core/internal/pkg/errs"
	"github.com/seanflannery10/core/internal/pkg/helpers"
	"github.com/seanflannery10/core/internal/pkg/validator"
	"github.com/seanflannery10/core/internal/services"
)

type createTokenAccessPayload struct {
	TokenPlaintext string `json:"token"`
}

func (p *createTokenAccessPayload) Bind(_ *http.Request) error {
	v := validator.New()

	data.ValidateTokenPlaintext(v, p.TokenPlaintext)

	if v.HasErrors() {
		return validator.NewValidationError(v.Errors)
	}

	return nil
}

func CreateTokenAccessHandler(env services.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p := &createTokenAccessPayload{}

		if helpers.CheckAndBind(w, r, p) {
			return
		}

		user, err := env.Queries.GetUserFromTokenHelper(r.Context(), p.TokenPlaintext, data.ScopeRefresh)
		if err != nil {
			_ = render.Render(w, r, errs.ErrServerError(err))
			return
		}

		token, err := env.Queries.CreateTokenHelper(r.Context(), user.ID, time.Hour, data.ScopeAccess)
		if err != nil {
			_ = render.Render(w, r, errs.ErrServerError(err))
			return
		}

		render.Status(r, http.StatusCreated)

		helpers.RenderAndCheck(w, r, &token)
	}
}
