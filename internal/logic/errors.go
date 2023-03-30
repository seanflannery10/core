package logic

import "github.com/go-faster/errors"

var (
	ErrActivationRequired   = errors.New("user account must be activated")
	ErrEditConflict         = errors.New("unable to update the record due to an edit conflict")
	ErrEmailNotFound        = errors.New("no matching email address found")
	ErrInvalidAccessToken   = errors.New("invalid access token")
	ErrInvalidCredentials   = errors.New("invalid authentication credentials")
	ErrInvalidToken         = errors.New("invalid or missing token")
	ErrMessageNotFound      = errors.New("no matching message found")
	ErrReusedRefreshToken   = errors.New("reused refresh token")
	ErrUserAlreadyActivated = errors.New("user has already been activated")
	ErrUserExists           = errors.New("a user with this email address already exists")
)
