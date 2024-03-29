package handler_test

import (
	"context"
	"regexp"
	"testing"
	"time"

	"github.com/go-faster/errors"
	"github.com/seanflannery10/core/internal/generated/api"
	"github.com/seanflannery10/core/internal/server/logic"
	"github.com/seanflannery10/core/internal/shared/utils"
	"github.com/stretchr/testify/assert"
)

const (
	tokenLength  = 26
	invalidToken = "token value is not valid"
)

func TestNewActivationToken_Success(t *testing.T) {
	request := &api.UserEmailRequest{
		Email: "unactivated@test.com",
	}

	response, err := newTestHandler(t).NewActivationToken(context.Background(), request)
	if err != nil {
		t.Fatalf(unexpectedError, err)
	}

	assert.Equal(t, logic.ScopeActivation, response.Scope)
	assert.Equal(t, tokenLength, len(response.Token))
	assert.IsType(t, time.Time{}, response.Expiry)

	if matches := regexp.MustCompile(`^([A-Za-z0-9+/]{4})*([A-Za-z0-9+/]{3}=|[A-Za-z0-9+/]{2}==)?$`).MatchString(response.Token); matches {
		t.Fatal(invalidToken)
	}
}

func TestNewActivationToken_NotFound(t *testing.T) {
	request := &api.UserEmailRequest{
		Email: "notfound@test.com",
	}

	response, err := newTestHandler(t).NewActivationToken(context.Background(), request)
	if !errors.Is(err, logic.ErrEmailNotFound) {
		t.Fatalf(unexpectedError, err)
	}

	if response != nil {
		t.Error(unexpectedResponse)
	}
}

func TestNewPasswordResetToken_Success(t *testing.T) {
	request := &api.UserEmailRequest{
		Email: "activated@test.com",
	}

	response, err := newTestHandler(t).NewPasswordResetToken(context.Background(), request)
	if err != nil {
		t.Fatalf(unexpectedError, err)
	}

	assert.Equal(t, logic.ScopePasswordReset, response.Scope)
	assert.Equal(t, tokenLength, len(response.Token))
	assert.IsType(t, time.Time{}, response.Expiry)

	if matches := regexp.MustCompile(`^([A-Za-z0-9+/]{4})*([A-Za-z0-9+/]{3}=|[A-Za-z0-9+/]{2}==)?$`).MatchString(response.Token); matches {
		t.Fatal(invalidToken)
	}
}

func TestNewPasswordResetToken_NotFound(t *testing.T) {
	request := &api.UserEmailRequest{
		Email: "notfound@test.com",
	}

	response, err := newTestHandler(t).NewPasswordResetToken(context.Background(), request)
	if !errors.Is(err, logic.ErrEmailNotFound) {
		t.Fatalf(unexpectedError, err)
	}

	if response != nil {
		t.Error(unexpectedResponse)
	}
}

func TestNewRefreshToken_Success(t *testing.T) {
	request := &api.UserLoginRequest{
		Email:    "activated@test.com",
		Password: "testtest",
	}

	tokenResponseHeaders, err := newTestHandler(t).NewRefreshToken(context.Background(), request)
	if err != nil {
		t.Fatalf(unexpectedError, err)
	}

	response := tokenResponseHeaders.Response
	cookie := tokenResponseHeaders.SetCookie.Value

	assert.Equal(t, logic.ScopeAccess, response.Scope)
	assert.Equal(t, tokenLength, len(response.Token))
	assert.IsType(t, time.Time{}, response.Expiry)

	if matches := regexp.MustCompile(`^([A-Za-z0-9+/]{4})*([A-Za-z0-9+/]{3}=|[A-Za-z0-9+/]{2}==)?$`).MatchString(response.Token); matches {
		t.Fatal(invalidToken)
	}

	assert.NotEmpty(t, cookie)
}

func TestNewRefreshToken_NotFound(t *testing.T) {
	request := &api.UserLoginRequest{
		Email:    "notfound1@test.com",
		Password: "testtest",
	}

	response, err := newTestHandler(t).NewRefreshToken(context.Background(), request)
	if !errors.Is(err, logic.ErrInvalidCredentials) {
		t.Fatalf(unexpectedError, err)
	}

	if response != nil {
		t.Error(unexpectedResponse)
	}
}

func TestNewAccessToken_Success(t *testing.T) {
	ctx := utils.ContextSetCookieValue(context.Background(), "AYNSWD44H2JKZAOHE3HF2BUK34")

	tokenResponseHeaders, err := newTestHandler(t).NewAccessToken(ctx)
	if err != nil {
		t.Fatalf(unexpectedError, err)
	}

	response := tokenResponseHeaders.Response
	cookie := tokenResponseHeaders.SetCookie.Value

	assert.Equal(t, logic.ScopeAccess, response.Scope)
	assert.Equal(t, tokenLength, len(response.Token))
	assert.IsType(t, time.Time{}, response.Expiry)

	if matches := regexp.MustCompile(`^([A-Za-z0-9+/]{4})*([A-Za-z0-9+/]{3}=|[A-Za-z0-9+/]{2}==)?$`).MatchString(response.Token); matches {
		t.Fatal(invalidToken)
	}

	assert.NotEmpty(t, cookie)
}

func TestNewAccessToken_NotFound(t *testing.T) {
	ctx := utils.ContextSetCookieValue(context.Background(), "NONE")

	response, err := newTestHandler(t).NewAccessToken(ctx)
	if !errors.Is(err, logic.ErrInvalidToken) {
		t.Fatalf(unexpectedError, err)
	}

	if response != nil {
		t.Error(unexpectedResponse)
	}
}
