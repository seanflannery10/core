package service

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/seanflannery10/core/internal/api"
	"github.com/seanflannery10/core/internal/data"
	"github.com/segmentio/asm/base64"
)

func newCookie(name, value string, secret []byte) (http.Cookie, error) {
	cookie := http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		MaxAge:   cookieTTL,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}

	block, err := aes.NewCipher(secret)
	if err != nil {
		return http.Cookie{}, fmt.Errorf("failed new cipher: %w", err)
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return http.Cookie{}, fmt.Errorf("failed new gcm: %w", err)
	}

	nonce := make([]byte, aesGCM.NonceSize())

	_, err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return http.Cookie{}, fmt.Errorf("failed read full: %w", err)
	}

	plaintext := fmt.Sprintf("%s:%s", cookie.Name, cookie.Value)

	encryptedValue := aesGCM.Seal(nonce, nonce, []byte(plaintext), nil)

	cookie.Value = base64.URLEncoding.EncodeToString(encryptedValue)

	if len(cookie.String()) > cookieMaxSize {
		return http.Cookie{}, fmt.Errorf("failed length check: %w", errValueTooLong)
	}

	return cookie, nil
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
