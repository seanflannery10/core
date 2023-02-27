package responses

import "net/http"

type StringResponsePayload struct {
	Message string `json:"message"`
}

func (p StringResponsePayload) Render(_ http.ResponseWriter, _ *http.Request) error {
	return nil
}
