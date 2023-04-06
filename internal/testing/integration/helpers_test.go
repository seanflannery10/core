package integration_test

import (
	"context"
	"encoding/hex"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/seanflannery10/core/internal/data"
	"github.com/seanflannery10/core/internal/handler"
	"github.com/seanflannery10/core/internal/shared/mailer"
	"github.com/seanflannery10/core/internal/shared/utils"
	"golang.org/x/exp/slog"
)

const (
	secretKey            = "ff2636f4a5abf829042c96d38caa8007427773980fddab20fd7c43d93dc186ca" //nolint:gosec
	connString           = "postgres://postgres:test@localhost:5433/test?sslmode=disable"
	page                 = 1
	pageSize             = 20
	testMessage          = "First!"
	testMessageEdit      = "Edit!"
	testMessageID        = 1
	testUserID           = 1
	testMessageIDMissing = 500
	testVersion          = 1
	testVersionEdit      = 2
	unexpectedError      = "unexpected error: %v"
	unexpectedResponse   = "unexpected response"
)

func newTestHandler() *handler.Handler {
	secret, err := hex.DecodeString(secretKey)
	if err != nil {
		slog.Error("filed decode string", err)
		os.Exit(1) //nolint:revive
	}

	dbpool, err := pgxpool.New(context.Background(), connString)
	if err != nil {
		slog.Error("filed new pool", err)
		os.Exit(1) //nolint:revive
	}

	mail, err := mailer.New(mailer.SMTP{Host: "localhost", Port: 2525})
	if err != nil {
		slog.Error("filed new mailer", err)
		os.Exit(1) //nolint:revive
	}

	return &handler.Handler{
		Queries: data.New(dbpool),
		Mailer:  mail,
		Secret:  secret,
	}
}

func ctxWithTestUser() context.Context {
	return utils.ContextSetUser(context.Background(), &data.User{
		ID:        testUserID,
		Name:      "test",
		Email:     "test@test.com",
		Activated: true,
	})
}
