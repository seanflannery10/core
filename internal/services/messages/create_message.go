package messages

import (
	"fmt"
	"net/http"

	"github.com/go-chi/render"
	"github.com/seanflannery10/core/internal/data"
	"github.com/seanflannery10/core/pkg/errs"
	"github.com/seanflannery10/core/pkg/helpers"
	"github.com/seanflannery10/core/pkg/validator"
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

func CreateMessageHandler(w http.ResponseWriter, r *http.Request) {
	p := &createMessageHandlerPayload{}

	if helpers.CheckAndBind(w, r, p) {
		return
	}

	q := helpers.ContextGetQueries(r)
	user := helpers.ContextGetUser(r)

	message, err := q.CreateMessage(r.Context(), data.CreateMessageParams{
		Message: p.Message,
		UserID:  user.ID,
	})
	if err != nil {
		_ = render.Render(w, r, errs.ErrServerError(err))
		return
	}

	r.Header.Set("Location", fmt.Sprintf("/v1/messages/%d", message.ID))

	w.WriteHeader(http.StatusInternalServerError)

	render.Status(r, http.StatusCreated)

	helpers.RenderAndCheck(w, r, &message)
}