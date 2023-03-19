package tokens

import (
	"crypto/sha256"
	"errors"
	"net/http"

	"github.com/go-chi/render"
	"github.com/jackc/pgx/v5"
	"github.com/seanflannery10/core/internal/data"
	"github.com/seanflannery10/core/internal/pkg/cookies"
	"github.com/seanflannery10/core/internal/pkg/errs"
	"github.com/seanflannery10/core/internal/pkg/helpers"
	"github.com/seanflannery10/core/internal/pkg/middleware"
	"github.com/seanflannery10/core/internal/services"
)

// @Summary	create access token using a refresh token
// @ID			create-token-access
// @Produce	json
// @Success	200	{object}	data.TokenFull
// @Router		/tokens/access  [post]
func CreateTokenAccessHandler(env *services.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		refreshTokenPlaintext, err := cookies.ReadEncrypted(r, cookieRefreshToken, env.Config.Secret)
		if err != nil {
			switch {
			case errors.Is(err, http.ErrNoCookie):
				_ = render.Render(w, r, errs.ErrCookieNotFound())
			case errors.Is(err, cookies.ErrInvalidValue):
				_ = render.Render(w, r, errs.ErrInvalidCookie())
			default:
				_ = render.Render(w, r, errs.ErrServerError(err))
			}

			return
		}

		user, err := env.Queries.GetUserFromTokenHelper(r.Context(), refreshTokenPlaintext, data.ScopeRefresh)
		if err != nil {
			switch {
			case errors.Is(err, pgx.ErrNoRows):
				_ = render.Render(w, r, errs.ErrInvalidToken())
			default:
				_ = render.Render(w, r, errs.ErrServerError(err))
			}

			return
		}

		tokenHash := sha256.Sum256([]byte(refreshTokenPlaintext))

		badToken, err := env.Queries.CheckToken(r.Context(), data.CheckTokenParams{
			Hash:   tokenHash[:],
			UserID: user.ID,
			Scope:  data.ScopeRefresh,
		})
		if err != nil {
			_ = render.Render(w, r, errs.ErrServerError(err))
			return
		}

		if badToken {
			err = env.Queries.DeleteTokens(r.Context(), data.DeleteTokensParams{
				Scope:  data.ScopeRefresh,
				UserID: user.ID,
			})
			if err != nil {
				_ = render.Render(w, r, errs.ErrServerError(err))
				return
			}

			_ = render.Render(w, r, errs.ErrReusedRefreshToken())
			return
		}

		err = env.Queries.DeactivateToken(r.Context(), data.DeactivateTokenParams{
			Scope:  data.ScopeRefresh,
			Hash:   tokenHash[:],
			UserID: user.ID,
		})
		if err != nil {
			_ = render.Render(w, r, errs.ErrServerError(err))
			return
		}

		env.User = &user

		w, accessToken, err := createRefreshAndAccessTokens(w, r, env)
		if err != nil {
			_ = render.Render(w, r, errs.ErrServerError(err))
			return
		}

		r = middleware.LogUser(r, &user)

		render.Status(r, http.StatusCreated)

		helpers.RenderAndCheck(w, r, &accessToken)
	}
}
