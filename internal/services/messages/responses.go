package messages

import (
	"net/http"

	"github.com/seanflannery10/core/internal/data"
	"github.com/seanflannery10/core/pkg/pagination"
)

type messagesResponsePayload struct {
	Messages []data.Message
	Metadata pagination.Metadata
}

func (p messagesResponsePayload) Render(_ http.ResponseWriter, _ *http.Request) error {
	return nil
}

type messageResponsePayload struct {
	Message data.Message
}

func (p messageResponsePayload) Render(_ http.ResponseWriter, _ *http.Request) error {
	return nil
}

type stringResponsePayload struct {
	Message string `json:"message"`
}

func (p stringResponsePayload) Render(_ http.ResponseWriter, _ *http.Request) error {
	return nil
}
