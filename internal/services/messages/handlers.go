package messages

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/render"
	"github.com/jackc/pgx/v5"
	"github.com/seanflannery10/core/internal/data"
	"github.com/seanflannery10/core/pkg/errs"
	"github.com/seanflannery10/core/pkg/helpers"
	"github.com/seanflannery10/core/pkg/responses"
	"github.com/seanflannery10/core/pkg/validator"
)

func createMessageHandler(w http.ResponseWriter, r *http.Request) {
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

	// TODO Fix Headers
	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/messages/%d", message.ID))

	render.Status(r, http.StatusCreated)

	helpers.RenderAndCheck(w, r, &message)
}

func showMessageHandler(w http.ResponseWriter, r *http.Request) {
	p := &showMessagePayload{}

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

func updateMessageHandler(w http.ResponseWriter, r *http.Request) {
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

func deleteMessageHandler(w http.ResponseWriter, r *http.Request) {
	p := &deleteMessagePayload{}

	if helpers.CheckAndBind(w, r, p) {
		return
	}

	q := helpers.ContextGetQueries(r)

	err := q.DeleteMessage(r.Context(), p.ID)
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

	helpers.RenderAndCheck(w, r, responses.NewStringResponsePayload("message successfully deleted"))
}

func listUserMessagesHandler(w http.ResponseWriter, r *http.Request) {
	p := &listUserMessagesPayload{}
	v := validator.New()

	if helpers.CheckAndBind(w, r, p) {
		return
	}

	q := helpers.ContextGetQueries(r)
	user := helpers.ContextGetUser(r)

	messages, err := q.GetUserMessages(r.Context(), data.GetUserMessagesParams{
		UserID: user.ID,
		Offset: p.Pagination.Offset(),
		Limit:  p.Pagination.Limit(),
	})
	if err != nil {
		_ = render.Render(w, r, errs.ErrServerError(err))
		return
	}

	count, err := q.GetUserMessageCount(r.Context(), user.ID)
	if err != nil {
		_ = render.Render(w, r, errs.ErrServerError(err))
		return
	}

	metadata := p.Pagination.CalculateMetadata(count, v)

	if v.HasErrors() {
		_ = render.Render(w, r, errs.ErrFailedValidation(v.Errors))
		return
	}

	render.Status(r, http.StatusCreated)

	helpers.RenderAndCheck(w, r, &messagesResponsePayload{messages, metadata})
}
