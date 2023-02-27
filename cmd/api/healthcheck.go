package main

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/seanflannery10/core/pkg/errs"
	"github.com/seanflannery10/core/pkg/helpers"
)

func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	env := map[string]any{
		"status": "available",
		"system_info": map[string]string{
			"version": helpers.GetVersion(),
		},
	}

	err := helpers.WriteJSON(w, http.StatusOK, env)
	if err != nil {
		_ = render.Render(w, r, errs.ErrServerError(err))
	}
}
