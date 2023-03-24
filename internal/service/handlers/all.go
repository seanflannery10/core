package handlers

import (
	"errors"
	"time"

	"github.com/seanflannery10/core/internal/shared/mailer"

	"github.com/seanflannery10/core/internal/data"
)

var (
	errNotActivated = errors.New("user not activated")

	errReusedRefreshToken = errors.New("reused refresh token")
	errAlreadyExists      = errors.New("already exists")
	errNotFound           = errors.New("not found")
	errInvalidCredentials = errors.New("invalid credentials")
)

const (
	lenthRandom = 16

	cookieRefreshToken = "core_refreshtoken"
	ttlCookie          = 7 * 24 * 60 * 60

	ttlAccessToken        = time.Hour
	ttlAcitvationToken    = 3 * 24 * time.Hour
	ttlPasswordResetToken = 45 * time.Minute
	ttlRefreshToken       = 7 * 24 * time.Hour
)

type Service struct {
	mailer  mailer.Mailer
	queries data.Queries
	secret  []byte
}
