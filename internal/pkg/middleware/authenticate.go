package middleware

import (
	"crypto/sha256"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/render"
	"github.com/jackc/pgx/v5"
	"github.com/seanflannery10/core/internal/data"
	"github.com/seanflannery10/core/internal/pkg/errs"
	"github.com/seanflannery10/core/internal/pkg/validator"
	"github.com/seanflannery10/core/internal/services"
)

func Authenticate(env *services.Env) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("Vary", "Authorization")

			authorizationHeader := r.Header.Get("Authorization")

			if authorizationHeader == "" {
				env.User = data.AnonymousUser
				next.ServeHTTP(w, r)
				return
			}

			headerParts := strings.Split(authorizationHeader, " ")
			if len(headerParts) != 2 || headerParts[0] != "Bearer" {
				_ = render.Render(w, r, errs.ErrInvalidAccessToken())
				return
			}

			token := headerParts[1]

			v := validator.New()

			data.ValidateTokenPlaintext(v, token)

			if v.HasErrors() {
				_ = render.Render(w, r, errs.ErrInvalidAccessToken())
				return
			}

			tokenHash := sha256.Sum256([]byte(token))

			user, err := env.Queries.GetUserFromToken(r.Context(), data.GetUserFromTokenParams{
				Hash:   tokenHash[:],
				Scope:  data.ScopeAccess,
				Expiry: time.Now(),
			})
			if err != nil {
				switch {
				case errors.Is(err, pgx.ErrNoRows):
					_ = render.Render(w, r, errs.ErrInvalidAccessToken())
				default:
					_ = render.Render(w, r, errs.ErrServerError(err))
				}

				return
			}

			env.User = &user

			r = LogUser(r, &user)

			next.ServeHTTP(w, r)
		})
	}
}
