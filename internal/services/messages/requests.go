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
	v       *validator.Validator
}

func (p *createMessageHandlerPayload) Bind(_ *http.Request) error {
	data.ValidateMessage(p.v, p.Message)

	if p.v.HasErrors() {
		return validator.ErrValidation
	}

	return nil
}

type showMessagePayload struct {
	ID int64
	v  *validator.Validator
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
	v       *validator.Validator
}

func (p *updateMessagePayload) Bind(r *http.Request) error {
	data.ValidateMessage(p.v, p.Message)

	if p.v.HasErrors() {
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
	v  *validator.Validator
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
	v *validator.Validator
}

func (p *listUserMessagesPayload) Bind(r *http.Request) error {
	p.Pagination = pagination.New(r, p.v)

	pagination.ValidatePagination(p.v, p.Pagination)

	if p.v.HasErrors() {
		return validator.ErrValidation
	}

	return nil
}
