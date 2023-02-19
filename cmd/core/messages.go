package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/jackc/pgx/v5"
	"github.com/seanflannery10/core/internal/data"
	"github.com/seanflannery10/core/internal/helpers"
	"github.com/seanflannery10/core/internal/httperrors"
	"github.com/seanflannery10/core/internal/pagination"
	"github.com/seanflannery10/core/internal/validator"
)

func (app *application) createMessageHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Message string `json:"message"`
	}

	err := helpers.ReadJSON(w, r, &input)
	if err != nil {
		httperrors.BadRequest(w, r, err)
		return
	}

	user := helpers.ContextGetUser(r)

	params := data.CreateMessageParams{
		Message: input.Message,
		UserID:  user.ID,
	}

	v := validator.New()

	if data.ValidateMessage(v, params.Message); v.HasErrors() {
		httperrors.FailedValidation(w, r, v)
		return
	}

	message, err := app.queries.CreateMessage(r.Context(), params)
	if err != nil {
		httperrors.ServerError(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/messages/%d", message.ID))

	err = helpers.WriteJSONWithHeaders(w, http.StatusCreated, map[string]any{"message": message}, headers)
	if err != nil {
		httperrors.ServerError(w, r, err)
	}
}

func (app *application) showMessageHandler(w http.ResponseWriter, r *http.Request) {
	id, err := helpers.ReadIDParam(r)
	if err != nil {
		httperrors.NotFound(w, r)
		return
	}

	message, err := app.queries.GetMessage(r.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			httperrors.NotFound(w, r)
		default:
			httperrors.ServerError(w, r, err)
		}

		return
	}

	err = helpers.WriteJSON(w, http.StatusOK, map[string]any{"message": message})
	if err != nil {
		httperrors.ServerError(w, r, err)
	}
}

func (app *application) updateMessageHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Message string `json:"message"`
	}

	err := helpers.ReadJSON(w, r, &input)
	if err != nil {
		httperrors.BadRequest(w, r, err)
		return
	}

	id, err := helpers.ReadIDParam(r)
	if err != nil {
		httperrors.NotFound(w, r)
		return
	}

	params := data.UpdateMessageParams{Message: input.Message, ID: id}

	v := validator.New()

	if data.ValidateMessage(v, params.Message); v.HasErrors() {
		httperrors.FailedValidation(w, r, v)
		return
	}

	message, err := app.queries.UpdateMessage(r.Context(), params)
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

	err = helpers.WriteJSON(w, http.StatusOK, map[string]any{"message": message})
	if err != nil {
		httperrors.ServerError(w, r, err)
	}
}

func (app *application) deleteMessageHandler(w http.ResponseWriter, r *http.Request) {
	id, err := helpers.ReadIDParam(r)
	if err != nil {
		httperrors.NotFound(w, r)
		return
	}

	err = app.queries.DeleteMessage(r.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			httperrors.NotFound(w, r)
		default:
			httperrors.ServerError(w, r, err)
		}

		return
	}

	err = helpers.WriteJSON(w, http.StatusOK, map[string]any{"message": "message successfully deleted"})
	if err != nil {
		httperrors.ServerError(w, r, err)
	}
}

func (app *application) listUserMessagesHandler(w http.ResponseWriter, r *http.Request) {
	v := validator.New()

	pgn := pagination.New(r, v)

	if pagination.ValidatePagination(v, pgn); v.HasErrors() {
		httperrors.FailedValidation(w, r, v)
		return
	}

	user := helpers.ContextGetUser(r)

	params := data.GetUserMessagesParams{
		UserID: user.ID,
		Offset: pgn.Offset(),
		Limit:  pgn.Limit(),
	}

	messages, err := app.queries.GetUserMessages(r.Context(), params)
	if err != nil {
		httperrors.ServerError(w, r, err)
		return
	}

	count, err := app.queries.GetUserMessageCount(r.Context(), user.ID)
	if err != nil {
		httperrors.ServerError(w, r, err)
		return
	}

	metadata := pgn.CalculateMetadata(count)

	err = helpers.WriteJSON(w, http.StatusOK, map[string]any{"messages": messages, "metadata": metadata})
	if err != nil {
		httperrors.ServerError(w, r, err)
	}
}
