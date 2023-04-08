package handler_test

import (
	"context"
	"testing"

	"github.com/go-faster/errors"
	"github.com/seanflannery10/core/internal/generated/api"
	"github.com/seanflannery10/core/internal/server/logic"
	"github.com/stretchr/testify/assert"
)

func TestActivateUser_Success(t *testing.T) {
	request := &api.TokenRequest{
		Token: "HJUKX2HGBVUJJ2R2RVGFB4RZ3I",
	}

	response, err := newTestHandler(t).ActivateUser(context.Background(), request)
	if err != nil {
		t.Fatalf(unexpectedError, err)
	}

	assert.Equal(t, "tobeactivated", response.Name)
	assert.Equal(t, "tobeactivated@test.com", response.Email)
}

func TestActivateUser_NotFound(t *testing.T) {
	request := &api.TokenRequest{
		Token: "NOTFOUND",
	}

	response, err := newTestHandler(t).ActivateUser(context.Background(), request)
	if !errors.Is(err, logic.ErrInvalidToken) {
		t.Fatalf(unexpectedError, err)
	}

	if response != nil {
		t.Error(unexpectedResponse)
	}
}

func TestNewUser_Success(t *testing.T) {
	request := &api.UserRequest{
		Name:     "newtest",
		Email:    "newtest@test.com",
		Password: "testtest",
	}

	response, err := newTestHandler(t).NewUser(context.Background(), request)
	if err != nil {
		t.Fatalf(unexpectedError, err)
	}

	assert.Equal(t, "newtest", response.Name)
	assert.Equal(t, "newtest@test.com", response.Email)
}

func TestNewUser_UnprocessableEntity(t *testing.T) {
	request := &api.UserRequest{
		Name:     "testexists",
		Email:    "activated@test.com",
		Password: "testtest",
	}

	response, err := newTestHandler(t).NewUser(context.Background(), request)
	if !errors.Is(err, logic.ErrUserExists) {
		t.Fatalf(unexpectedError, err)
	}

	if response != nil {
		t.Error(unexpectedResponse)
	}
}

func TestUpdateUserPassword_Success(t *testing.T) {
	request := &api.UpdateUserPasswordRequest{
		Password: "newtestpass",
		Token:    "DH332OIAAI3JHJ3VHN5IPIZAB4",
	}

	response, err := newTestHandler(t).UpdateUserPassword(context.Background(), request)
	if err != nil {
		t.Fatalf(unexpectedError, err)
	}

	assert.Equal(t, "password updated", response.Message)
}

func TestUpdateUserPassword_NotFound(t *testing.T) {
	request := &api.UpdateUserPasswordRequest{
		Password: "newtestpass",
		Token:    "NOTFOUND",
	}

	response, err := newTestHandler(t).UpdateUserPassword(context.Background(), request)
	if !errors.Is(err, logic.ErrInvalidToken) {
		t.Fatalf(unexpectedError, err)
	}

	if response != nil {
		t.Error(unexpectedResponse)
	}
}
