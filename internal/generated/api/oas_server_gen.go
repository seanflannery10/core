// Code generated by ogen, DO NOT EDIT.

package api

import (
	"context"
)

// Handler handles operations described by OpenAPI v3 specification.
type Handler interface {
	// ActivateUser implements ActivateUser operation.
	//
	// PATCH /v1/users/activate
	ActivateUser(ctx context.Context, req *TokenRequest) (*UserResponse, error)
	// DeleteMessage implements DeleteMessage operation.
	//
	// DELETE /v1/messages/{id}
	DeleteMessage(ctx context.Context, params DeleteMessageParams) (*AcceptanceResponse, error)
	// GetMessage implements GetMessage operation.
	//
	// GET /v1/messages/{id}
	GetMessage(ctx context.Context, params GetMessageParams) (*MessageResponse, error)
	// GetUserMessages implements GetUserMessages operation.
	//
	// GET /v1/messages
	GetUserMessages(ctx context.Context, params GetUserMessagesParams) (*MessagesResponse, error)
	// NewAccessToken implements NewAccessToken operation.
	//
	// POST /v1/tokens/access
	NewAccessToken(ctx context.Context) (*TokenResponseHeaders, error)
	// NewActivationToken implements NewActivationToken operation.
	//
	// POST /v1/tokens/activation
	NewActivationToken(ctx context.Context, req *UserEmailRequest) (*TokenResponse, error)
	// NewMessage implements NewMessage operation.
	//
	// POST /v1/messages
	NewMessage(ctx context.Context, req *MessageRequest) (*MessageResponse, error)
	// NewPasswordResetToken implements NewPasswordResetToken operation.
	//
	// POST /v1/tokens/password-reset
	NewPasswordResetToken(ctx context.Context, req *UserEmailRequest) (*TokenResponse, error)
	// NewRefreshToken implements NewRefreshToken operation.
	//
	// POST /v1/tokens/refresh
	NewRefreshToken(ctx context.Context, req *UserLoginRequest) (*TokenResponseHeaders, error)
	// NewUser implements NewUser operation.
	//
	// POST /v1/users/register
	NewUser(ctx context.Context, req *UserRequest) (*UserResponse, error)
	// UpdateMessage implements UpdateMessage operation.
	//
	// PUT /v1/messages/{id}
	UpdateMessage(ctx context.Context, req *MessageRequest, params UpdateMessageParams) (*MessageResponse, error)
	// UpdateUserPassword implements UpdateUserPassword operation.
	//
	// PATCH /v1/users/update-password
	UpdateUserPassword(ctx context.Context, req *UpdateUserPasswordRequest) (*AcceptanceResponse, error)
	// NewError creates *ErrorResponseStatusCode from error returned by handler.
	//
	// Used for common default response.
	NewError(ctx context.Context, err error) *ErrorResponseStatusCode
}

// Server implements http server based on OpenAPI v3 specification and
// calls Handler to handle requests.
type Server struct {
	h   Handler
	sec SecurityHandler
	baseServer
}

// NewServer creates new Server.
func NewServer(h Handler, sec SecurityHandler, opts ...ServerOption) (*Server, error) {
	s, err := newServerConfig(opts...).baseServer()
	if err != nil {
		return nil, err
	}
	return &Server{
		h:          h,
		sec:        sec,
		baseServer: s,
	}, nil
}