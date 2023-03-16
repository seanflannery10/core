package messages

import (
	"fmt"
	"net/http"

	"github.com/go-chi/render"
	"github.com/seanflannery10/core/internal/data"
	"github.com/seanflannery10/core/internal/pkg/errs"
	"github.com/seanflannery10/core/internal/pkg/helpers"
	"github.com/seanflannery10/core/internal/pkg/validator"
	"github.com/seanflannery10/core/internal/services"
)

type createMessageHandlerPayload struct {
	Message string `json:"message"`
}

func (p *createMessageHandlerPayload) Bind(_ *http.Request) error {
	v := validator.New()

	data.ValidateMessage(v, p.Message)

	if v.HasErrors() {
		return validator.NewValidationError(v.Errors)
	}

	return nil
}

func CreateMessageHandler(env *services.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p := &createMessageHandlerPayload{}

		if helpers.CheckAndBind(w, r, p) {
			return
		}

		user := helpers.ContextGetUser(r)

		message, err := env.Queries.CreateMessage(r.Context(), data.CreateMessageParams{
			Message: p.Message,
			UserID:  user.ID,
		})
		if err != nil {
			_ = render.Render(w, r, errs.ErrServerError(err))
			return
		}

		r.Header.Set("Location", fmt.Sprintf("/v1/messages/%d", message.ID))

		render.Status(r, http.StatusCreated)

		helpers.RenderAndCheck(w, r, &message)
	}
}
