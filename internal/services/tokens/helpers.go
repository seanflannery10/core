package tokens

import (
	"fmt"
	"net/http"
	"time"

	"github.com/seanflannery10/core/internal/data"
	"github.com/seanflannery10/core/internal/pkg/cookies"
	"github.com/seanflannery10/core/internal/services"
)

const (
	cookieRefreshToken    = "core_refreshtoken"
	ttlAccessToken        = time.Hour
	ttlAcitvationToken    = 3 * 24 * time.Hour
	ttlCookie             = 7 * 24 * 60 * 60
	ttlPasswordResetToken = 45 * time.Minute
	ttlRefreshToken       = 7 * 24 * time.Hour
)

func createCookie(w http.ResponseWriter, name, value string, secret []byte) (http.ResponseWriter, error) {
	tokenCookie := &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		MaxAge:   ttlCookie,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}

	err := cookies.WriteEncrypted(w, tokenCookie, secret)
	if err != nil {
		return nil, fmt.Errorf("failed write encrypted: %w", err)
	}

	return w, nil
}

func createRefreshAndAccessTokens(w http.ResponseWriter, r *http.Request, env *services.Env) (http.ResponseWriter, data.TokenFull, error) { //nolint:revive
	refreshToken, err := env.Queries.CreateTokenHelper(r.Context(), env.User.ID, ttlRefreshToken, data.ScopeRefresh)
	if err != nil {
		return nil, data.TokenFull{}, fmt.Errorf("failed create token helper: %w", err)
	}

	w, err = createCookie(w, cookieRefreshToken, refreshToken.Plaintext, env.Config.Secret)
	if err != nil {
		return nil, data.TokenFull{}, fmt.Errorf("failed create cookie: %w", err)
	}

	accessToken, err := env.Queries.CreateTokenHelper(r.Context(), env.User.ID, ttlAccessToken, data.ScopeAccess)
	if err != nil {
		return nil, data.TokenFull{}, fmt.Errorf("failed create token helper: %w", err)
	}

	return w, accessToken, nil
}
