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

	if errors.Is(err, errUserNotActivated) {
		code = http.StatusForbidden
	}

	w.Header().Set("WWW-Authenticate", "Bearer")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	errStr := strings.ReplaceAll(errors.Unwrap(err).Error(), "\"", "'")

	e := jx.GetEncoder()
	e.ObjStart()
	e.FieldStart("error")
	e.StrEscape(errStr)
	e.ObjEnd()

	_, _ = w.Write(e.Bytes())
}
