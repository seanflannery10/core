package messages

import (
	"errors"
	"net/http"

	"github.com/go-chi/render"
	"github.com/jackc/pgx/v5"
	"github.com/seanflannery10/core/internal/pkg/errs"
	"github.com/seanflannery10/core/internal/pkg/helpers"
	"github.com/seanflannery10/core/internal/services"
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

func GetMessageHandler(env services.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p := &getMessagePayload{}

		if helpers.CheckAndBind(w, r, p) {
			return
		}

		message, err := env.Queries.GetMessage(r.Context(), p.ID)
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
}
