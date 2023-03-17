package services

import (
	"github.com/seanflannery10/core/internal/data"
	"github.com/seanflannery10/core/internal/pkg/mailer"
	oteltrace "go.opentelemetry.io/otel/trace"
)

type (
	Env struct {
		Queries *data.Queries
		Mailer  mailer.Mailer
		Tracer  oteltrace.Tracer
		User    data.User
		Config  Config
	}

	Config struct {
		SMTP         mailer.SMTP
		Env          string `env:"ENV,default=dev"`
		OTelEndpoint string `env:"OTEL_EXPORTER_OTLP_ENDPOINT,default=api.honeycomb.io:443"`
		DatabaseURL  string `env:"DATABASE_URL,default=postgres://postgres:test@localhost:5432/test?sslmode=disable"`
		SecretKey    string `env:"SECRET_KEY"`
		Secret       []byte `env:"SECRET_KEY"`
		Port         int    `env:"PORT,default=4000"`
	}
)
