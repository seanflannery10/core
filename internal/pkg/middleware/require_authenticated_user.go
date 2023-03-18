package middleware

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/seanflannery10/core/internal/pkg/errs"
	"github.com/seanflannery10/core/internal/services"
)

func RequireAuthenticatedUser(env *services.Env) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if env.User.IsAnonymous() {
				_ = render.Render(w, r, errs.ErrAuthenticationRequired())
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
