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

func TestGetUserMessages_Empty(t *testing.T) {
	dbpool, err := pgxpool.New(context.Background(), connString)
	if err != nil {
		t.Fatalf(unexpectedError, err)
	}

	ctx := context.Background()
	q := data.New(dbpool)
	expectedResponse := &api.MessagesResponse{
		Messages: []api.MessageResponse{},
		Metadata: api.MessagesMetadataResponse{},
	}

	response, err := logic.GetUserMessages(ctx, q, page, pageSize, testUserID)
	if err != nil {
		t.Fatalf(unexpectedError, err)
	}

	assert.Equal(t, expectedResponse, response)
}

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
		Message: testMessage,
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

	response, err := logic.GetMessage(ctx, q, testMessageIDMissing, testUserID)
	if !errors.Is(err, logic.ErrMessageNotFound) {
		t.Fatalf(unexpectedError, err)
	}

	if response != nil {
		t.Error(unexpectedResponse)
	}
}

func TestUpdateMessage_Success(t *testing.T) {
	dbpool, err := pgxpool.New(context.Background(), connString)
	if err != nil {
		t.Fatalf(unexpectedError, err)
	}

	ctx := context.Background()
	q := data.New(dbpool)
	expectedResponse := &api.MessageResponse{
		ID:      testMessageID,
		Message: testMessageEdit,
		Version: testVersionEdit,
	}

	response, err := logic.UpdateMessage(ctx, q, testMessageEdit, testMessageID, testUserID)
	if err != nil {
		t.Fatalf(unexpectedError, err)
	}

	assert.Equal(t, expectedResponse, response)
}

func TestUpdateMessage_NotFound(t *testing.T) {
	dbpool, err := pgxpool.New(context.Background(), connString)
	if err != nil {
		t.Fatalf(unexpectedError, err)
	}

	ctx := context.Background()
	q := data.New(dbpool)

	response, err := logic.UpdateMessage(ctx, q, testMessageEdit, testMessageIDMissing, testUserID)
	if !errors.Is(err, logic.ErrMessageNotFound) {
		t.Fatalf(unexpectedError, err)
	}

	if response != nil {
		t.Error(unexpectedResponse)
	}
}

func TestGetUserMessages_WithMessage(t *testing.T) {
	dbpool, err := pgxpool.New(context.Background(), connString)
	if err != nil {
		t.Fatalf(unexpectedError, err)
	}

	ctx := context.Background()
	q := data.New(dbpool)
	expectedResponse := &api.MessagesResponse{
		Messages: []api.MessageResponse{
			{ID: testMessageID, Message: testMessageEdit, Version: testVersionEdit},
		},
		Metadata: api.MessagesMetadataResponse{
			CurrentPage:  page,
			FirstPage:    page,
			LastPage:     page,
			PageSize:     pageSize,
			TotalRecords: page,
		},
	}

	response, err := logic.GetUserMessages(ctx, q, page, pageSize, testUserID)
	if err != nil {
		t.Fatalf(unexpectedError, err)
	}

	assert.Equal(t, expectedResponse, response)
}

func TestDeleteMessage_Success(t *testing.T) {
	dbpool, err := pgxpool.New(context.Background(), connString)
	if err != nil {
		t.Fatalf(unexpectedError, err)
	}

	ctx := context.Background()
	q := data.New(dbpool)
	expectedResponse := &api.AcceptanceResponse{Message: "message deleted"}

	response, err := logic.DeleteMessage(ctx, q, testMessageID, testUserID)
	if err != nil {
		t.Fatalf(unexpectedError, err)
	}

	assert.Equal(t, expectedResponse, response)
}

func TestDeleteMessage_NotFound(t *testing.T) {
	dbpool, err := pgxpool.New(context.Background(), connString)
	if err != nil {
		t.Fatalf(unexpectedError, err)
	}

	ctx := context.Background()
	q := data.New(dbpool)

	response, err := logic.DeleteMessage(ctx, q, testMessageIDMissing, testUserID)
	if !errors.Is(err, logic.ErrMessageNotFound) {
		t.Fatalf(unexpectedError, err)
	}

	if response != nil {
		t.Error(unexpectedResponse)
	}
}
