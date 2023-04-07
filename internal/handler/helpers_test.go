package handler_test

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
	connString = "postgres://postgres:test@localhost:5433/test?sslmode=disable"
	secretKey  = "ff2636f4a5abf829042c96d38caa8007427773980fddab20fd7c43d93dc186ca" //nolint:gosec
	exitCode   = 1
)

func newTestHandler() *handler.Handler {
	secret, err := hex.DecodeString(secretKey)
	if err != nil {
		slog.Error("filed decode string", err)
		os.Exit(exitCode) //nolint:revive
	}

	dbpool, err := pgxpool.New(context.Background(), connString)
	if err != nil {
		slog.Error("filed new pool", err)
		os.Exit(exitCode) //nolint:revive
	}

	mail, err := mailer.New(mailer.SMTP{Host: "localhost", Port: 2525, Sender: "Test <no-reply@testdomain.com>", Test: true})
	if err != nil {
		slog.Error("filed new mailer", err)
		os.Exit(exitCode) //nolint:revive
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
		Email:     "activated@test.com",
		Activated: true,
	})
}
