package tokens

import (
	"net/http"

	"github.com/seanflannery10/core/internal/data"
)

type tokenResponsePayload struct {
	Token data.TokenFull
}

func (p tokenResponsePayload) Render(_ http.ResponseWriter, _ *http.Request) error {
	return nil
}
