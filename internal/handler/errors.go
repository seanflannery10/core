package handler

import (
	"context"
	"net/http"
	"strings"

	"github.com/go-faster/errors"
	"github.com/go-faster/jx"
	"github.com/ogen-go/ogen/ogenerrors"
	"github.com/seanflannery10/core/internal/api"
	"github.com/seanflannery10/core/internal/data"
	"github.com/seanflannery10/core/internal/logic"
	"github.com/seanflannery10/core/internal/shared/mailer"
	"github.com/seanflannery10/core/internal/shared/pagination"
	"golang.org/x/exp/slog"
)

type Handler struct {
	Mailer  mailer.Mailer
	Queries *data.Queries
	Secret  []byte
}

func (s *Handler) NewError(_ context.Context, err error) *api.ErrorResponseStatusCode {
	var (
		code       int
		errMessage = errors.Unwrap(err).Error()

		activationRequired   = errors.Is(err, logic.ErrActivationRequired)
		editConflict         = errors.Is(err, logic.ErrEditConflict)
		emailNotFound        = errors.Is(err, logic.ErrEmailNotFound)
		invalidCredentials   = errors.Is(err, logic.ErrInvalidCredentials)
		invalidToken         = errors.Is(err, logic.ErrInvalidToken)
		messageNotFound      = errors.Is(err, logic.ErrMessageNotFound)
		pageValueToHigh      = errors.Is(err, pagination.ErrPageValueToHigh)
		reusedRefreshToken   = errors.Is(err, logic.ErrReusedRefreshToken)
		userAlreadyActivated = errors.Is(err, logic.ErrUserAlreadyActivated)
		userExists           = errors.Is(err, logic.ErrUserExists)
	)

	switch {
	case invalidCredentials, reusedRefreshToken:
		code = http.StatusUnauthorized
	case emailNotFound, messageNotFound:
		code = http.StatusNotFound
	case editConflict:
		code = http.StatusConflict
	case activationRequired, invalidToken, pageValueToHigh, userAlreadyActivated, userExists:
		code = http.StatusUnprocessableEntity
	default:
		slog.Error("server error", "error", err)

		code = http.StatusInternalServerError
		errMessage = "the server encountered a problem and could not process your request"
	}

	return &api.ErrorResponseStatusCode{StatusCode: code, Response: api.ErrorResponse{Error: errMessage}}
}

func ErrorHandler(_ context.Context, w http.ResponseWriter, _ *http.Request, err error) {
	code := ogenerrors.ErrorCode(err)
	errMessage := strings.ReplaceAll(errors.Unwrap(err).Error(), "\"", "'")

	if errors.Is(err, ogenerrors.ErrSecurityRequirementIsNotSatisfied) {
		errMessage = "missing security token or cookie"
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	e := jx.GetEncoder()
	e.ObjStart()
	e.FieldStart("error")
	e.StrEscape(errMessage)
	e.ObjEnd()

	_, _ = w.Write(e.Bytes())
}
