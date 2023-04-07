package handler_test

import (
	"context"
	"encoding/hex"
	"testing"

	"github.com/go-faster/errors"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/seanflannery10/core/internal/data"
	"github.com/seanflannery10/core/internal/handler"
	"github.com/seanflannery10/core/internal/shared/mailer"
	"github.com/seanflannery10/core/internal/shared/utils"
)

const (
	connString = "postgres://postgres:test@localhost:5433/test?sslmode=disable"
	secretKey  = "ff2636f4a5abf829042c96d38caa8007427773980fddab20fd7c43d93dc186ca" //nolint:gosec
)

func newTestHandler(t *testing.T) *handler.Handler {
	t.Helper()

	secret, err := hex.DecodeString(secretKey)
	if err != nil {
		t.Fatalf(unexpectedError, errors.Wrap(err, "filed decode string"))
	}

	dbpool, err := pgxpool.New(context.Background(), connString)
	if err != nil {
		t.Fatalf(unexpectedError, errors.Wrap(err, "filed new pool"))
	}

	mail, err := mailer.New(mailer.SMTP{Host: "localhost", Port: 2525, Sender: "Test <no-reply@testdomain.com>", Test: true})
	if err != nil {
		t.Fatalf(unexpectedError, errors.Wrap(err, "filed new mailer"))
	}

	return &handler.Handler{
		Queries: data.New(dbpool),
		Mailer:  mail,
		Secret:  secret,
	}
}

func ctxWithTestUser(t *testing.T) context.Context {
	t.Helper()

	return utils.ContextSetUser(context.Background(), &data.User{
		ID:        testUserID,
		Name:      "test",
		Email:     "activated@test.com",
		Activated: true,
	})
}
