package logic

import (
	"time"

	"github.com/go-faster/errors"
)

const (
	ttlAccessToken        = time.Hour
	ttlActivationToken    = 3 * 24 * time.Hour
	ttlPasswordResetToken = 45 * time.Minute
	ttlRefreshToken       = 7 * 24 * time.Hour
)

var (
	errAlreadyExists      = errors.New("already exists")
	errInvalidCredentials = errors.New("invalid credentials")
	errNotActivated       = errors.New("not activated")
	errNotFound           = errors.New("not found")
	errReusedRefreshToken = errors.New("reused refresh token")
)
