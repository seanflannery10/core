package responses

import "net/http"

type StringResponsePayload struct {
	Message string `json:"message"`
}

func NewStringResponsePayload(m string) *StringResponsePayload {
	return &StringResponsePayload{Message: m}
}

func (p *StringResponsePayload) Render(_ http.ResponseWriter, _ *http.Request) error {
	return nil
}
