package main

import (
	"context"
	"net/http"
	"strings"

	"github.com/go-faster/errors"
	"github.com/go-faster/jx"
	"github.com/ogen-go/ogen/ogenerrors"
)

func (app *application) ErrorHandler(_ context.Context, w http.ResponseWriter, _ *http.Request, err error) {
	code := ogenerrors.ErrorCode(err)
	errMessage := strings.ReplaceAll(errors.Unwrap(err).Error(), "\"", "'")

	switch {
	case errors.Is(err, errUserNotActivated):
		code = http.StatusForbidden
	case errors.Is(err, ogenerrors.ErrSecurityRequirementIsNotSatisfied):
		errMessage = "missing security token or cookie"
	}

	// w.Header().Set("WWW-Authenticate", "Bearer")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	e := jx.GetEncoder()
	e.ObjStart()
	e.FieldStart("error")
	e.StrEscape(errMessage)
	e.ObjEnd()

	_, _ = w.Write(e.Bytes())
}
