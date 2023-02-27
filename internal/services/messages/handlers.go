package messages

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/render"
	"github.com/jackc/pgx/v5"
	"github.com/seanflannery10/core/internal/data"
	"github.com/seanflannery10/core/pkg/helpers"
	"github.com/seanflannery10/core/pkg/httperrors"
	"github.com/seanflannery10/core/pkg/validator"
	"golang.org/x/exp/slog"
)

func createMessageHandler(w http.ResponseWriter, r *http.Request) {
	p := &createMessageHandlerPayload{v: validator.New()}

	err := render.Bind(r, p)
	if err != nil {
		helpers.CheckBindErr(w, r, p.v, err)
		return
	}

	q := helpers.ContextGetQueries(r)
	user := helpers.ContextGetUser(r)

	message, err := q.CreateMessage(r.Context(), data.CreateMessageParams{
		Message: p.Message,
		UserID:  user.ID,
	})
	if err != nil {
		httperrors.ServerError(w, r, err)
		return
	}

	// TODO Fix Headers
	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/messages/%d", message.ID))

	render.Status(r, http.StatusCreated)

	err = render.Render(w, r, messageResponsePayload{message})
	if err != nil {
		slog.Error("render error", err)
	}
}

func showMessageHandler(w http.ResponseWriter, r *http.Request) {
	p := &showMessagePayload{v: validator.New()}

	err := render.Bind(r, p)
	if err != nil {
		helpers.CheckBindErr(w, r, p.v, err)
		return
	}

	q := helpers.ContextGetQueries(r)

	message, err := q.GetMessage(r.Context(), p.ID)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			httperrors.NotFound(w, r)
		default:
			httperrors.ServerError(w, r, err)
		}

		return
	}

	render.Status(r, http.StatusCreated)

	err = render.Render(w, r, messageResponsePayload{message})
	if err != nil {
		slog.Error("render error", err)
	}
}

func updateMessageHandler(w http.ResponseWriter, r *http.Request) {
	p := &updateMessagePayload{v: validator.New()}

	err := render.Bind(r, p)
	if err != nil {
		helpers.CheckBindErr(w, r, p.v, err)
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
			httperrors.NotFound(w, r)
		default:
			httperrors.ServerError(w, r, err)
		}

		return
	}

	if r.Header.Get("X-Expected-Version") != "" {
		if strconv.FormatInt(int64(message.Version), 32) != r.Header.Get("X-Expected-Version") {
			httperrors.EditConflict(w, r)
			return
		}
	}

	render.Status(r, http.StatusCreated)

	err = render.Render(w, r, messageResponsePayload{message})
	if err != nil {
		slog.Error("render error", err)
	}
}

func deleteMessageHandler(w http.ResponseWriter, r *http.Request) {
	p := &deleteMessagePayload{v: validator.New()}

	err := render.Bind(r, p)
	if err != nil {
		helpers.CheckBindErr(w, r, p.v, err)
		return
	}

	q := helpers.ContextGetQueries(r)

	err = q.DeleteMessage(r.Context(), p.ID)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			httperrors.NotFound(w, r)
		default:
			httperrors.ServerError(w, r, err)
		}

		return
	}

	render.Status(r, http.StatusCreated)

	err = render.Render(w, r, stringResponsePayload{"message successfully deleted"})
	if err != nil {
		slog.Error("render error", err)
	}
}

func listUserMessagesHandler(w http.ResponseWriter, r *http.Request) {
	p := &listUserMessagesPayload{v: validator.New()}

	err := render.Bind(r, p)
	if err != nil {
		helpers.CheckBindErr(w, r, p.v, err)
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
		httperrors.ServerError(w, r, err)
		return
	}

	count, err := q.GetUserMessageCount(r.Context(), user.ID)
	if err != nil {
		httperrors.ServerError(w, r, err)
		return
	}

	metadata := p.Pagination.CalculateMetadata(count, p.v)

	if p.v.HasErrors() {
		httperrors.FailedValidation(w, r, p.v)
		return
	}

	render.Status(r, http.StatusCreated)

	err = render.Render(w, r, messagesResponsePayload{messages, metadata})
	if err != nil {
		slog.Error("render error", err)
	}
}
