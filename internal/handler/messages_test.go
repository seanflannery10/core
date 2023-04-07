package handler_test

import (
	"testing"

	"github.com/go-faster/errors"
	"github.com/seanflannery10/core/internal/api"
	"github.com/seanflannery10/core/internal/logic"
	"github.com/stretchr/testify/assert"
)

const (
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

func TestGetUserMessages_SuccessEmpty(t *testing.T) {
	params := api.GetUserMessagesParams{
		Page:     api.OptInt32{Value: page, Set: true},
		PageSize: api.OptInt32{Value: pageSize, Set: true},
	}

	expected := &api.MessagesResponse{Messages: []api.MessageResponse{}, Metadata: api.MessagesMetadataResponse{}}

	response, err := newTestHandler(t).GetUserMessages(ctxWithTestUser(t), params)
	if err != nil {
		t.Fatalf(unexpectedError, err)
	}

	assert.Equal(t, expected, response)
}

func TestNewMessage_Success(t *testing.T) {
	request := &api.MessageRequest{Message: testMessage}

	expected := &api.MessageResponse{ID: testMessageID, Message: testMessage, Version: testVersion}

	response, err := newTestHandler(t).NewMessage(ctxWithTestUser(t), request)
	if err != nil {
		t.Fatalf(unexpectedError, err)
	}

	assert.Equal(t, expected, response)
}

func TestGetMessage_Success(t *testing.T) {
	params := api.GetMessageParams{ID: testMessageID}

	expected := &api.MessageResponse{ID: testMessageID, Message: testMessage, Version: testVersion}

	response, err := newTestHandler(t).GetMessage(ctxWithTestUser(t), params)
	if err != nil {
		t.Fatalf(unexpectedError, err)
	}

	assert.Equal(t, expected, response)
}

func TestGetMessage_NotFound(t *testing.T) {
	params := api.GetMessageParams{ID: testMessageIDMissing}

	response, err := newTestHandler(t).GetMessage(ctxWithTestUser(t), params)
	if !errors.Is(err, logic.ErrMessageNotFound) {
		t.Fatalf(unexpectedError, err)
	}

	if response != nil {
		t.Error(unexpectedResponse)
	}
}

func TestUpdateMessage_Success(t *testing.T) {
	request := &api.MessageRequest{Message: testMessageEdit}
	params := api.UpdateMessageParams{ID: testMessageID}

	expected := &api.MessageResponse{ID: testMessageID, Message: testMessageEdit, Version: testVersionEdit}

	response, err := newTestHandler(t).UpdateMessage(ctxWithTestUser(t), request, params)
	if err != nil {
		t.Fatalf(unexpectedError, err)
	}

	assert.Equal(t, expected, response)
}

func TestUpdateMessage_NotFound(t *testing.T) {
	request := &api.MessageRequest{Message: testMessageEdit}
	params := api.UpdateMessageParams{ID: testMessageIDMissing}

	response, err := newTestHandler(t).UpdateMessage(ctxWithTestUser(t), request, params)
	if !errors.Is(err, logic.ErrMessageNotFound) {
		t.Fatalf(unexpectedError, err)
	}

	if response != nil {
		t.Error(unexpectedResponse)
	}
}

func TestGetUserMessages_SuccessWithMessage(t *testing.T) {
	params := api.GetUserMessagesParams{
		Page:     api.OptInt32{Value: page, Set: true},
		PageSize: api.OptInt32{Value: pageSize, Set: true},
	}

	expected := &api.MessagesResponse{
		Messages: []api.MessageResponse{{ID: testMessageID, Message: testMessageEdit, Version: testVersionEdit}},
		Metadata: api.MessagesMetadataResponse{CurrentPage: page, FirstPage: page, LastPage: page, PageSize: pageSize, TotalRecords: page},
	}

	response, err := newTestHandler(t).GetUserMessages(ctxWithTestUser(t), params)
	if err != nil {
		t.Fatalf(unexpectedError, err)
	}

	assert.Equal(t, expected, response)
}

func TestDeleteMessage_Success(t *testing.T) {
	params := api.DeleteMessageParams{ID: testMessageID}

	expected := &api.AcceptanceResponse{Message: "message deleted"}

	response, err := newTestHandler(t).DeleteMessage(ctxWithTestUser(t), params)
	if err != nil {
		t.Fatalf(unexpectedError, err)
	}

	assert.Equal(t, expected, response)
}

func TestDeleteMessage_NotFound(t *testing.T) {
	params := api.DeleteMessageParams{ID: testMessageIDMissing}

	response, err := newTestHandler(t).DeleteMessage(ctxWithTestUser(t), params)
	if !errors.Is(err, logic.ErrMessageNotFound) {
		t.Fatalf(unexpectedError, err)
	}

	if response != nil {
		t.Error(unexpectedResponse)
	}
}
