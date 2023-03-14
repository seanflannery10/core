package tokens

import (
	"encoding/hex"
	"errors"
	"net/http"
	"os"

	"github.com/go-chi/render"
	"github.com/seanflannery10/core/internal/data"
	"github.com/seanflannery10/core/internal/pkg/cookies"
	"github.com/seanflannery10/core/internal/pkg/errs"
	"github.com/seanflannery10/core/internal/pkg/helpers"
	"github.com/seanflannery10/core/internal/services"
)

func CreateTokenAccessHandler(env services.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		secret, err := hex.DecodeString(os.Getenv("SECRET_KEY"))
		if err != nil {
			_ = render.Render(w, r, errs.ErrServerError(err))
		}

		refreshTokenPlaintext, err := cookies.ReadEncrypted(r, "core_refreshtoken", secret)
		if err != nil {
			switch {
			case errors.Is(err, http.ErrNoCookie):
				_ = render.Render(w, r, &errs.ErrResponse{
					AppCode:        0,
					HTTPStatusCode: http.StatusBadRequest,
					Message:        "cookie not found",
				})
			case errors.Is(err, cookies.ErrInvalidValue):
				_ = render.Render(w, r, &errs.ErrResponse{
					AppCode:        0,
					HTTPStatusCode: http.StatusBadRequest,
					Message:        "invalid cookie",
				})
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

		w, accessToken, err := createRefreshAndAccessTokens(w, r, env.Queries, user.ID)
		if err != nil {
			_ = render.Render(w, r, errs.ErrServerError(err))
			return
		}

		render.Status(r, http.StatusCreated)

		helpers.RenderAndCheck(w, r, &accessToken)
	}
}
