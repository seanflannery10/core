package service

import (
	"time"

	"github.com/go-faster/errors"
	"github.com/seanflannery10/core/internal/data"
	"github.com/seanflannery10/core/internal/shared/mailer"
)

var (
	errAlreadyExists      = errors.New("already exists")
	errInvalidCredentials = errors.New("invalid credentials")
	errNotActivated       = errors.New("not activated")
	errNotFound           = errors.New("not found")
	errReusedRefreshToken = errors.New("reused refresh token")
)

const (
	cookieRefreshToken = "core_refreshtoken"
	cookieTTL          = 7 * 24 * 60 * 60

	ttlAccessToken        = time.Hour
	ttlAcitvationToken    = 3 * 24 * time.Hour
	ttlPasswordResetToken = 45 * time.Minute
	ttlRefreshToken       = 7 * 24 * time.Hour
)

type Handler struct {
	Mailer  mailer.Mailer
	Queries data.Queries
	Secret  []byte
}
