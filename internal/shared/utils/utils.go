package utils

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
	"runtime/debug"
	"time"

	"github.com/go-faster/errors"
	"github.com/ogen-go/ogen/middleware"
	"github.com/seanflannery10/core/internal/data"
	"github.com/seanflannery10/core/internal/oas"
	"github.com/segmentio/asm/base64"
)

var errValueTooLong = errors.New("cookie value too long")

const (
	cookieMaxSize  = 4096
	lenthRandom    = 16
	userContextKey = contextKey("user")
)

type contextKey string

func ContextSetUser(req *middleware.Request, user *data.User) *middleware.Request {
	ctx := context.WithValue(req.Context, userContextKey, *user)

	req.Raw = req.Raw.WithContext(ctx)

	return req
}

func ContextGetUser(ctx context.Context) data.User {
	user, ok := ctx.Value(userContextKey).(data.User)
	if !ok {
		panic("missing user value in request context")
	}

	return user
}

func NewCookie(name, value string, ttl int, secret []byte) (oas.OptString, error) {
	cookie := http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		MaxAge:   ttl,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}

	block, err := aes.NewCipher(secret)
	if err != nil {
		return oas.OptString{}, fmt.Errorf("failed new cipher: %w", err)
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return oas.OptString{}, fmt.Errorf("failed new gcm: %w", err)
	}

	nonce := make([]byte, aesGCM.NonceSize())

	_, err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return oas.OptString{}, fmt.Errorf("failed read full: %w", err)
	}

	plaintext := fmt.Sprintf("%s:%s", cookie.Name, cookie.Value)

	encryptedValue := aesGCM.Seal(nonce, nonce, []byte(plaintext), nil)

	cookie.Value = base64.URLEncoding.EncodeToString(encryptedValue)

	if len(cookie.String()) > cookieMaxSize {
		return oas.OptString{}, fmt.Errorf("failed length check: %w", errValueTooLong)
	}

	optString := oas.OptString{Value: cookie.Value, Set: true}

	return optString, nil
}

func NewToken(ctx context.Context, q data.Queries, ttl time.Duration, scope string, userID int64) (oas.TokenResponse, error) {
	randomBytes := make([]byte, lenthRandom)

	_, err := rand.Read(randomBytes)
	if err != nil {
		return oas.TokenResponse{}, fmt.Errorf("failed read rand: %w", err)
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
		return oas.TokenResponse{}, fmt.Errorf("failed create token: %w", err)
	}

	tokenPlaintext := oas.TokenResponse{
		Plaintext: plaintext,
		Expiry:    oas.OptDateTime{Value: token.Expiry, Set: true},
		Scope:     oas.OptString{Value: token.Scope, Set: true},
	}

	return tokenPlaintext, nil
}

func GetVersion() string {
	var (
		revision string
		modified bool
	)

	bi, ok := debug.ReadBuildInfo()
	if ok {
		for _, s := range bi.Settings {
			switch s.Key {
			case "vcs.revision":
				revision = s.Value
			case "vcs.modified":
				if s.Value == "true" {
					modified = true
				}
			}
		}
	}

	if revision == "" {
		return "unavailable"
	}

	if modified {
		return fmt.Sprintf("%s-dirty", revision)
	}

	return revision
}
