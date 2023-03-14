package tokens

import (
	"errors"
	"net/http"

	"github.com/go-chi/render"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
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

		if sessionID != "" {
			err = env.Queries.DeleteSessionTokensForUser(r.Context(), data.DeleteSessionTokensForUserParams{
				UserID:  user.ID,
				Session: pgtype.Text{String: sessionID, Valid: true},
			})
			if err != nil {
				_ = render.Render(w, r, errs.ErrServerError(err))
				return
			}
		}

		newSessionID := uuid.New().String()

		w, err = createCookie(w, cookieSessionID, newSessionID, env.Config.Secret)
		if err != nil {
			_ = render.Render(w, r, errs.ErrServerError(err))
			return
		}

		w, accessToken, err := createRefreshAndAccessTokens(w, r, env, user.ID, newSessionID)
		if err != nil {
			_ = render.Render(w, r, errs.ErrServerError(err))
			return
		}

		render.Status(r, http.StatusCreated)

		helpers.RenderAndCheck(w, r, &accessToken)
	}
}
