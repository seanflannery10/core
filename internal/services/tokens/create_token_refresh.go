package tokens

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
	"encoding/hex"
	"errors"
	"net/http"
	"os"
	"time"

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

		sessionID, err := cookies.ReadEncrypted(r, "core_sessionid", env.Config.Secret)
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

		err = env.Queries.DeleteSessionTokensForUser(r.Context(), data.DeleteSessionTokensForUserParams{
			UserID:  user.ID,
			Session: pgtype.Text{String: sessionID, Valid: true},
		})
		if err != nil {
			_ = render.Render(w, r, errs.ErrServerError(err))
			return
		}

		w, newSessionID, err := createSessionCookie(w)
		if err != nil {
			_ = render.Render(w, r, errs.ErrServerError(err))
			return
		}

		w, accessToken, err := createRefreshAndAccessTokens(w, r, env.Queries, user.ID, newSessionID)
		if err != nil {
			_ = render.Render(w, r, errs.ErrServerError(err))
			return
		}

		render.Status(r, http.StatusCreated)

		helpers.RenderAndCheck(w, r, &accessToken)
	}
}

func createRefreshCookie(w http.ResponseWriter, plaintextToken string) (http.ResponseWriter, error) {
	tokenCookie := http.Cookie{
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

	err = cookies.WriteEncrypted(w, tokenCookie, secret)
	if err != nil {
		return nil, err
	}

	return w, nil
}

func createSessionCookie(w http.ResponseWriter) (http.ResponseWriter, string, error) {
	uid := uuid.New().String()

	sessionCookie := http.Cookie{
		Name:     "core_sessionid",
		Value:    uid,
		Path:     "/",
		MaxAge:   7 * 24 * 60 * 60,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}

	secret, err := hex.DecodeString(os.Getenv("SECRET_KEY"))
	if err != nil {
		return nil, "", err
	}

	err = cookies.WriteEncrypted(w, sessionCookie, secret)
	if err != nil {
		return nil, "", err
	}

	return w, uid, nil
}

func createRefreshAndAccessTokens(w http.ResponseWriter, r *http.Request, q *data.Queries, uid int64, sid string) (http.ResponseWriter, data.TokenFull, error) { //nolint:lll
	randomBytes := make([]byte, 16)

	_, err := rand.Read(randomBytes)
	if err != nil {
		return nil, data.TokenFull{}, err
	}

	plaintext := base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(randomBytes)
	hash := sha256.Sum256([]byte(plaintext))

	token, err := q.CreateRefreshToken(r.Context(), data.CreateRefreshTokenParams{
		Hash:    hash[:],
		UserID:  uid,
		Expiry:  pgtype.Timestamptz{Time: time.Now().Add(7 * 24 * time.Hour), Valid: true},
		Session: pgtype.Text{String: sid, Valid: true},
	})
	if err != nil {
		return nil, data.TokenFull{}, err
	}

	refreshToken := data.TokenFull{
		Plaintext: plaintext,
		Hash:      token.Hash,
		UserID:    token.UserID,
		Expiry:    token.Expiry,
		Scope:     token.Scope,
	}

	w, err = createRefreshCookie(w, refreshToken.Plaintext)
	if err != nil {
		return nil, data.TokenFull{}, err
	}

	accessToken, err := q.CreateTokenHelper(r.Context(), uid, time.Hour, data.ScopeAccess)
	if err != nil {
		return nil, data.TokenFull{}, err
	}

	return w, accessToken, nil
}
