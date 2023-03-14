package services

import (
	"github.com/seanflannery10/core/internal/data"
	"github.com/seanflannery10/core/internal/pkg/mailer"
	oteltrace "go.opentelemetry.io/otel/trace"
)

type Env struct {
	Queries *data.Queries
	Mailer  mailer.Mailer
	Tracer  oteltrace.Tracer
}
