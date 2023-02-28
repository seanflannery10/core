package messages

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/seanflannery10/core/internal/data"
	"github.com/seanflannery10/core/pkg/errs"
	"github.com/seanflannery10/core/pkg/helpers"
	"github.com/seanflannery10/core/pkg/pagination"
	"github.com/seanflannery10/core/pkg/validator"
)

type listMessagesUserPayload struct {
	pagination.Pagination
}

func (p *listMessagesUserPayload) Bind(r *http.Request) error {
	v := validator.New()

	p.Pagination = pagination.New(r, v)

	pagination.ValidatePagination(v, p.Pagination)

	if v.HasErrors() {
		return validator.NewValidationError(v.Errors)
	}

	return nil
}

type messagesResponsePayload struct {
	Messages []data.Message
	Metadata pagination.Metadata
}

func (p messagesResponsePayload) Render(_ http.ResponseWriter, _ *http.Request) error {
	return nil
}

func ListMessagesUserHandler(w http.ResponseWriter, r *http.Request) {
	p := &listMessagesUserPayload{}
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
