package integration_test

import (
	"context"
	"testing"

	"github.com/go-faster/errors"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/seanflannery10/core/internal/api"
	"github.com/seanflannery10/core/internal/data"
	"github.com/seanflannery10/core/internal/logic"
	"github.com/stretchr/testify/assert"
)

const (
	missingMessageID = 500
	testMessage      = "First!"
	connString       = "postgres://postgres:test@localhost:5433/test?sslmode=disable"
	testMessageID    = 1
	testUserID       = 1
	testVersion      = 1
	unexpectedError  = "unexpected error: %v"
)

func TestNewMessage_Success(t *testing.T) {
	dbpool, err := pgxpool.New(context.Background(), connString)
	if err != nil {
		t.Fatalf(unexpectedError, err)
	}

	ctx := context.Background()
	q := data.New(dbpool)
	expectedResponse := &api.MessageResponse{
		ID:      testMessageID,
		Message: testMessage,
		Version: testVersion,
	}

	response, err := logic.NewMessage(ctx, q, testMessage, testUserID)
	if err != nil {
		t.Fatalf(unexpectedError, err)
	}

	assert.Equal(t, expectedResponse, response)
}

func TestGetMessage_Success(t *testing.T) {
	dbpool, err := pgxpool.New(context.Background(), connString)
	if err != nil {
		t.Fatalf(unexpectedError, err)
	}

	ctx := context.Background()
	q := data.New(dbpool)
	expectedResponse := &api.MessageResponse{
		ID:      testMessageID,
		Message: "First!",
		Version: testVersion,
	}

	response, err := logic.GetMessage(ctx, q, testMessageID, testUserID)
	if err != nil {
		t.Fatalf(unexpectedError, err)
	}

	assert.Equal(t, expectedResponse, response)
}

func TestGetMessage_NotFound(t *testing.T) {
	dbpool, err := pgxpool.New(context.Background(), connString)
	if err != nil {
		t.Fatalf(unexpectedError, err)
	}

	ctx := context.Background()
	q := data.New(dbpool)

	response, err := logic.GetMessage(ctx, q, missingMessageID, testUserID)
	if !errors.Is(err, logic.ErrMessageNotFound) {
		t.Fatalf(unexpectedError, err)
	}

	if response != nil {
		t.Error("unexpected response")
	}
}
