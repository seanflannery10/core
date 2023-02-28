package messages

import (
	"net/http"

	"github.com/seanflannery10/core/internal/data"
	"github.com/seanflannery10/core/pkg/helpers"
	"github.com/seanflannery10/core/pkg/pagination"
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

type showMessagePayload struct {
	ID int64
}

func (p *showMessagePayload) Bind(r *http.Request) error {
	id, err := helpers.ReadIDParam(r)
	if err != nil {
		return err
	}

	p.ID = id

	return nil
}

type updateMessagePayload struct {
	Message string `json:"message"`
	ID      int64
}

func (p *updateMessagePayload) Bind(r *http.Request) error {
	v := validator.New()

	data.ValidateMessage(v, p.Message)

	if v.HasErrors() {
		return validator.ErrValidation
	}

	id, err := helpers.ReadIDParam(r)
	if err != nil {
		return err
	}

	p.ID = id

	return nil
}

type deleteMessagePayload struct {
	ID int64
}

func (p *deleteMessagePayload) Bind(r *http.Request) error {
	id, err := helpers.ReadIDParam(r)
	if err != nil {
		return err
	}

	p.ID = id

	return nil
}

type listUserMessagesPayload struct {
	pagination.Pagination
}

func (p *listUserMessagesPayload) Bind(r *http.Request) error {
	v := validator.New()

	p.Pagination = pagination.New(r, v)

	pagination.ValidatePagination(v, p.Pagination)

	if v.HasErrors() {
		return validator.ErrValidation
	}

	return nil
}
