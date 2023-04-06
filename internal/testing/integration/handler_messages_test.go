package integration_test

import (
	"testing"

	"github.com/seanflannery10/core/internal/api"
	"github.com/stretchr/testify/assert"
)

func TestGetUserMessagesHandler_Success(t *testing.T) {
	params := api.GetUserMessagesParams{
		Page:     api.OptInt32{Value: page, Set: true},
		PageSize: api.OptInt32{Value: pageSize, Set: true},
	}

	expected := &api.MessagesResponse{
		Messages: []api.MessageResponse{},
		Metadata: api.MessagesMetadataResponse{},
	}

	response, err := newTestHandler().GetUserMessages(ctxWithTestUser(), params)
	if err != nil {
		t.Fatalf(unexpectedError, err)
	}

	assert.Equal(t, expected, response)
}
