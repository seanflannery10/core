package handler_test

import (
	"context"
	"regexp"
	"testing"
	"time"

	"github.com/go-faster/errors"
	"github.com/seanflannery10/core/internal/api"
	"github.com/seanflannery10/core/internal/logic"
	"github.com/stretchr/testify/assert"
)

const (
	tokenLength = 26
)

func TestNewActivationToken_Success(t *testing.T) {
	request := &api.UserEmailRequest{
		Email: "unactivated@test.com",
	}

	response, err := newTestHandler().NewActivationToken(context.Background(), request)
	if err != nil {
		t.Fatalf(unexpectedError, err)
	}

	assert.Equal(t, logic.ScopeActivation, response.Scope)
	assert.Equal(t, tokenLength, len(response.Token))
	assert.IsType(t, time.Time{}, response.Expiry)

	if matches := regexp.MustCompile(`^([A-Za-z0-9+/]{4})*([A-Za-z0-9+/]{3}=|[A-Za-z0-9+/]{2}==)?$`).MatchString(response.Token); matches {
		t.Fatal("token value is not valid")
	}
}

func TestNewActivationToken_NotFound(t *testing.T) {
	request := &api.UserEmailRequest{
		Email: "notfound@test.com",
	}

	response, err := newTestHandler().NewActivationToken(context.Background(), request)
	if !errors.Is(err, logic.ErrEmailNotFound) {
		t.Fatalf(unexpectedError, err)
	}

	if response != nil {
		t.Error(unexpectedResponse)
	}
}
