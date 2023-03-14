package tokens

import (
	"errors"
	"net/http"

	"github.com/go-chi/render"
	"github.com/seanflannery10/core/internal/data"
	"github.com/seanflannery10/core/internal/pkg/cookies"
	"github.com/seanflannery10/core/internal/pkg/errs"
	"github.com/seanflannery10/core/internal/pkg/helpers"
	"github.com/seanflannery10/core/internal/services"
)

func CreateTokenAccessHandler(env services.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		refreshTokenPlaintext, err := cookies.ReadEncrypted(r, cookieSessionID, env.Config.Secret)
		if err != nil {
			switch {
			case errors.Is(err, http.ErrNoCookie):
				_ = render.Render(w, r, errs.ErrCookieNotFound)
			case errors.Is(err, cookies.ErrInvalidValue):
				_ = render.Render(w, r, errs.ErrInvalidCookie)
			default:
				_ = render.Render(w, r, errs.ErrServerError(err))
			}

			return
		}

		user, err := env.Queries.GetUserFromTokenHelper(r.Context(), refreshTokenPlaintext, data.ScopeRefresh)
		if err != nil {
			_ = render.Render(w, r, errs.ErrServerError(err))
			return
		}

		sessionID, err := cookies.ReadEncrypted(r, cookieSessionID, env.Config.Secret)
		if err != nil {
			switch {
			case errors.Is(err, http.ErrNoCookie):
				_ = render.Render(w, r, errs.ErrCookieNotFound)
			case errors.Is(err, cookies.ErrInvalidValue):
				_ = render.Render(w, r, errs.ErrInvalidCookie)
			default:
				_ = render.Render(w, r, errs.ErrServerError(err))
			}

			return
		}

		// TODO Check if refresh token is current and unused

		w, accessToken, err := createRefreshAndAccessTokens(w, r, env, user.ID, sessionID)
		if err != nil {
			_ = render.Render(w, r, errs.ErrServerError(err))
			return
		}

		render.Status(r, http.StatusCreated)

		helpers.RenderAndCheck(w, r, &accessToken)
	}
}
