package main

import (
	"net/http"

	"github.com/seanflannery10/core/internal/helpers"
	"github.com/seanflannery10/core/internal/httperrors"
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
		httperrors.ServerError(w, r, err)
	}
}
