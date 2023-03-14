package tokens

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/seanflannery10/core/internal/data"
	"github.com/seanflannery10/core/internal/pkg/cookies"
	"github.com/seanflannery10/core/internal/services"
)

const (
	cookieRefreshToken = "core_refreshtoken"
	cookieSessionID    = "core_sessionid"
)

func createCookie(w http.ResponseWriter, name, value string, secret []byte) (http.ResponseWriter, error) {
	tokenCookie := http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		MaxAge:   7 * 24 * 60 * 60,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}

	err := cookies.WriteEncrypted(w, tokenCookie, secret)
	if err != nil {
		return nil, err
	}

	return w, nil
}

func createRefreshAndAccessTokens(w http.ResponseWriter, r *http.Request, env services.Env, uid int64, sid string) (http.ResponseWriter, data.TokenFull, error) { //nolint:lll
	randomBytes := make([]byte, 16)

	_, err := rand.Read(randomBytes)
	if err != nil {
		return nil, data.TokenFull{}, err
	}

	plaintext := base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(randomBytes)
	hash := sha256.Sum256([]byte(plaintext))

	_, err = env.Queries.CreateRefreshToken(r.Context(), data.CreateRefreshTokenParams{
		Hash:    hash[:],
		UserID:  uid,
		Expiry:  pgtype.Timestamptz{Time: time.Now().Add(7 * 24 * time.Hour), Valid: true},
		Session: pgtype.Text{String: sid, Valid: true},
	})
	if err != nil {
		return nil, data.TokenFull{}, err
	}

	w, err = createCookie(w, cookieRefreshToken, plaintext, env.Config.Secret)
	if err != nil {
		return nil, data.TokenFull{}, err
	}

	accessToken, err := env.Queries.CreateTokenHelper(r.Context(), uid, time.Hour, data.ScopeAccess)
	if err != nil {
		return nil, data.TokenFull{}, err
	}

	return w, accessToken, nil
}
