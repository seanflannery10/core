package messages

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/render"
	"github.com/jackc/pgx/v5"
	"github.com/seanflannery10/core/internal/data"
	"github.com/seanflannery10/core/pkg/errs"
	"github.com/seanflannery10/core/pkg/helpers"
	"github.com/seanflannery10/core/pkg/validator"
)

type updateMessagePayload struct {
	Message string `json:"message"`
	ID      int64
}

func (p *updateMessagePayload) Bind(r *http.Request) error {
	v := validator.New()

	data.ValidateMessage(v, p.Message)

	if v.HasErrors() {
		return validator.NewValidationError(v.Errors)
	}

	id, err := helpers.ReadIDParam(r)
	if err != nil {
		return err
	}

	p.ID = id

	return nil
}

func UpdateMessageHandler(w http.ResponseWriter, r *http.Request) {
	p := &updateMessagePayload{}

	if helpers.CheckAndBind(w, r, p) {
		return
	}

	q := helpers.ContextGetQueries(r)

	message, err := q.UpdateMessage(r.Context(), data.UpdateMessageParams{
		Message: p.Message,
		ID:      p.ID,
	})
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			_ = render.Render(w, r, errs.ErrNotFound)
		default:
			_ = render.Render(w, r, errs.ErrServerError(err))
		}

		return
	}

	if r.Header.Get("X-Expected-Version") != "" {
		if strconv.FormatInt(int64(message.Version), 32) != r.Header.Get("X-Expected-Version") {
			_ = render.Render(w, r, errs.ErrEditConflict)
			return
		}
	}

	render.Status(r, http.StatusCreated)

	helpers.RenderAndCheck(w, r, &message)
}
