package tokens

import (
	"encoding/hex"
	"errors"
	"net/http"
	"os"
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

		w, accessToken, err := createRefreshAndAccessTokens(w, r, env.Queries, user.ID)
		if err != nil {
			_ = render.Render(w, r, errs.ErrServerError(err))
			return
		}

		render.Status(r, http.StatusCreated)

		helpers.RenderAndCheck(w, r, &accessToken)
	}
}

func createRefreshCookie(w http.ResponseWriter, plaintextToken string) (http.ResponseWriter, error) {
	cookie := http.Cookie{
		Name:     "core_refreshtoken",
		Value:    plaintextToken,
		Path:     "/",
		MaxAge:   7 * 24 * 60 * 60,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}

	secret, err := hex.DecodeString(os.Getenv("SECRET_KEY"))
	if err != nil {
		return nil, err
	}

	err = cookies.WriteEncrypted(w, cookie, secret)
	if err != nil {
		return nil, err
	}

	return w, nil
}

func createRefreshAndAccessTokens(w http.ResponseWriter, r *http.Request, q *data.Queries, id int64) (http.ResponseWriter, data.TokenFull, error) { //nolint:lll
	refreshToken, err := q.CreateTokenHelper(r.Context(), id, 7*24*time.Hour, data.ScopeRefresh)
	if err != nil {
		return nil, data.TokenFull{}, err
	}

	w, err = createRefreshCookie(w, refreshToken.Plaintext)
	if err != nil {
		return nil, data.TokenFull{}, err
	}

	accessToken, err := q.CreateTokenHelper(r.Context(), id, time.Hour, data.ScopeAccess)
	if err != nil {
		return nil, data.TokenFull{}, err
	}

	return w, accessToken, nil
}
