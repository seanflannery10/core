package messages

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/render"
	"github.com/jackc/pgx/v5"
	"github.com/seanflannery10/core/internal/pkg/errs"
	"github.com/seanflannery10/core/internal/pkg/helpers"
	"github.com/seanflannery10/core/internal/pkg/responses"
	"github.com/seanflannery10/core/internal/services"
)

type deleteMessagePayload struct {
	ID int64 `json:"-"`
}

func (p *deleteMessagePayload) Bind(r *http.Request) error {
	id, err := helpers.ReadIDParam(r)
	if err != nil {
		return fmt.Errorf("failed read id: %w", err)
	}

	p.ID = id

	return nil
}

func DeleteMessageHandler(env *services.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p := &deleteMessagePayload{}

		if helpers.CheckAndBind(w, r, p) {
			return
		}

		err := env.Queries.DeleteMessage(r.Context(), p.ID)
		if err != nil {
			switch {
			case errors.Is(err, pgx.ErrNoRows):
				_ = render.Render(w, r, errs.ErrNotFound())
			default:
				_ = render.Render(w, r, errs.ErrServerError(err))
			}

			return
		}

		render.Status(r, http.StatusCreated)

		helpers.RenderAndCheck(w, r, responses.NewStringResponsePayload("message successfully deleted"))
	}
}
