package handler

import (
	"github.com/seanflannery10/core/internal/data"
	"github.com/seanflannery10/core/internal/oas"
	"github.com/seanflannery10/core/internal/shared/mailer"
)

const (
	cookieRefreshToken = "core_refreshtoken"
	cookieTTL          = 7 * 24 * 60 * 60
)

type Handler struct {
	oas.UnimplementedHandler

	Mailer  mailer.Mailer
	Queries data.Queries
	Secret  []byte
}
