package services

import (
	"github.com/seanflannery10/core/internal/data"
	"github.com/seanflannery10/core/pkg/mailer"
	"go.opentelemetry.io/otel/trace"
)

type Env struct {
	Queries   *data.Queries
	Mailer    mailer.Mailer
	StdTracer *trace.Tracer
	ErrTracer *trace.Tracer
}
