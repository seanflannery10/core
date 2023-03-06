package services

import (
	"github.com/seanflannery10/core/internal/data"
	"github.com/seanflannery10/core/pkg/mailer"
	"github.com/seanflannery10/core/pkg/telemetry"
)

type Env struct {
	Queries *data.Queries
	Mailer  mailer.Mailer
	Tracers telemetry.Tracers
}
