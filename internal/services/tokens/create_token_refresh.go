package tokens

import (
	"encoding/hex"
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/render"
	"github.com/jackc/pgx/v5"
	"github.com/seanflannery10/core/internal/data"
	"github.com/seanflannery10/core/internal/pkg/cookies"
	"github.com/seanflannery10/core/internal/pkg/errs"
	"github.com/seanflannery10/core/internal/pkg/helpers"
	"github.com/seanflannery10/core/internal/pkg/validator"
	"github.com/seanflannery10/core/internal/services"
)

type createTokenRefreshPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (p *createTokenRefreshPayload) Bind(_ *http.Request) error {
	v := validator.New()

	data.ValidateEmail(v, p.Email)
	data.ValidatePasswordPlaintext(v, p.Password)

	if v.HasErrors() {
		return validator.NewValidationError(v.Errors)
	}

	return nil
}

func CreateTokenRefreshHandler(env services.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p := &createTokenRefreshPayload{}

		if helpers.CheckAndBind(w, r, p) {
			return
		}

		user, err := env.Queries.GetUserFromEmail(r.Context(), p.Email)
		if err != nil {
			switch {
			case errors.Is(err, pgx.ErrNoRows):
				_ = render.Render(w, r, errs.ErrInvalidCredentials)
			default:
				_ = render.Render(w, r, errs.ErrServerError(err))
			}

			return
		}

		match, err := user.ComparePasswords(p.Password)
		if err != nil {
			_ = render.Render(w, r, errs.ErrServerError(err))
			return
		}

		if !match {
			_ = render.Render(w, r, errs.ErrInvalidCredentials)
			return
		}

		refreshToken, err := env.Queries.CreateTokenHelper(r.Context(), user.ID, 30*24*time.Hour, data.ScopeRefresh)
		if err != nil {
			_ = render.Render(w, r, errs.ErrServerError(err))
			return
		}

		cookie := http.Cookie{
			Name:     "core_refresh_token",
			Value:    refreshToken.Plaintext,
			Path:     "/",
			MaxAge:   int(30 * 24 * time.Hour),
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteLaxMode,
		}

		secret, err := hex.DecodeString(env.Config.SecretKey)
		if err != nil {
			_ = render.Render(w, r, errs.ErrServerError(err))
			return
		}

		err = cookies.WriteEncrypted(w, cookie, secret)
		if err != nil {
			_ = render.Render(w, r, errs.ErrServerError(err))
			return
		}

		scopeToken, err := env.Queries.CreateTokenHelper(r.Context(), user.ID, time.Hour, data.ScopeAccess)
		if err != nil {
			_ = render.Render(w, r, errs.ErrServerError(err))
			return
		}

		render.Status(r, http.StatusCreated)

		helpers.RenderAndCheck(w, r, &scopeToken)
	}
}
