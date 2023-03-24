package handlers

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
	"fmt"
	"net/http"
	"time"

	"github.com/seanflannery10/core/internal/api"
	"github.com/seanflannery10/core/internal/data"
	"github.com/seanflannery10/core/internal/shared/cookies"
)

func NewCookie(w http.ResponseWriter, name, value string, secret []byte) (http.ResponseWriter, error) {
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

func newToken(ctx context.Context, q data.Queries, ttl time.Duration, scope string, userID int64) (api.TokenResponse, error) {
	randomBytes := make([]byte, lenthRandom)

	_, err := rand.Read(randomBytes)
	if err != nil {
		return api.TokenResponse{}, fmt.Errorf("failed read rand: %w", err)
	}

	plaintext := base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(randomBytes)
	hash := sha256.Sum256([]byte(plaintext))

	token, err := q.CreateToken(ctx, data.CreateTokenParams{
		Hash:   hash[:],
		UserID: userID,
		Expiry: time.Now().Add(ttl),
		Scope:  scope,
	})
	if err != nil {
		return api.TokenResponse{}, fmt.Errorf("failed create token: %w", err)
	}

	tokenPlaintext := api.TokenResponse{
		Plaintext: plaintext,
		Expiry:    api.OptDateTime{Value: token.Expiry, Set: true},
		Scope:     api.OptString{Value: token.Scope, Set: true},
	}

	return tokenPlaintext, nil
}
