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
