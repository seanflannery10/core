package messages

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/seanflannery10/core/internal/data"
	"github.com/seanflannery10/core/pkg/errs"
	"github.com/seanflannery10/core/pkg/helpers"
	pagination "github.com/seanflannery10/core/pkg/pagination"
	"github.com/seanflannery10/core/pkg/validator"
)

type getMessagesUserPayload struct {
	pagination.Pagination `json:"-"`
}

func (p *getMessagesUserPayload) Bind(r *http.Request) error {
	p.Pagination.Validate()

	if p.Pagination.Validator.HasErrors() {
		return validator.NewValidationError(p.Pagination.Validator.Errors)
	}

	return nil
}

func GetMessagesUserHandler(w http.ResponseWriter, r *http.Request) {
	p := &getMessagesUserPayload{Pagination: pagination.New(r)}

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

	metadata := p.Pagination.CalculateMetadata(count)

	if p.Pagination.Validator.HasErrors() {
		_ = render.Render(w, r, errs.ErrFailedValidation(p.Pagination.Validator.Errors))
		return
	}

	render.Status(r, http.StatusCreated)

	helpers.RenderAndCheck(w, r, &messagesResponsePayload{Messages: messages, Metadata: metadata})
}

type messagesResponsePayload struct {
	Messages []data.Message      `json:"messages"`
	Metadata pagination.Metadata `json:"metadata"`
}

func (p messagesResponsePayload) Render(_ http.ResponseWriter, _ *http.Request) error {
	return nil
}
