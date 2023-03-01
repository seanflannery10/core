package messages

import (
	"errors"
	"net/http"

	"github.com/go-chi/render"
	"github.com/jackc/pgx/v5"
	"github.com/seanflannery10/core/pkg/errs"
	"github.com/seanflannery10/core/pkg/helpers"
)

type getMessagePayload struct {
	ID int64 `json:"-"`
}

func (p *getMessagePayload) Bind(r *http.Request) error {
	id, err := helpers.ReadIDParam(r)
	if err != nil {
		return err
	}

	p.ID = id

	return nil
}

func GetMessageHandler(w http.ResponseWriter, r *http.Request) {
	p := &getMessagePayload{}

	if helpers.CheckAndBind(w, r, p) {
		return
	}

	q := helpers.ContextGetQueries(r)

	message, err := q.GetMessage(r.Context(), p.ID)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			_ = render.Render(w, r, errs.ErrNotFound)
		default:
			_ = render.Render(w, r, errs.ErrServerError(err))
		}

		return
	}

	render.Status(r, http.StatusCreated)

	helpers.RenderAndCheck(w, r, &message)
}
