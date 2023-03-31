package handler

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"

	"github.com/go-faster/errors"
)

var errValueTooLong = errors.New("cookie value too long")

const (
	cookieMaxSize      = 4096
	cookieRefreshToken = "core_refresh_token"
	cookieTTL          = 7 * 24 * 60 * 60
)

func newCookie(name, value string, ttl int, secret []byte) (http.Cookie, error) {
	cookie := http.Cookie{Name: name, Value: value, Path: "/", MaxAge: ttl, HttpOnly: true, Secure: true, SameSite: http.SameSiteLaxMode}

	block, err := aes.NewCipher(secret)
	if err != nil {
		return http.Cookie{}, fmt.Errorf("failed new cipher: %w", err)
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return http.Cookie{}, fmt.Errorf("failed new gcm: %w", err)
	}

	nonce := make([]byte, aesGCM.NonceSize())

	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
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
