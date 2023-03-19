package data

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
	"fmt"
	"net/http"
	"time"

	"github.com/seanflannery10/core/internal/pkg/validator"
)

const (
	ScopeAccess        = "access"
	ScopeActivation    = "activation"
	ScopePasswordReset = "password-reset"
	ScopeRefresh       = "refresh"
	emptyString        = ""
	keyToken           = "token"
	lengthToken        = 26
	lenthRandom        = 16
)

type TokenFull struct {
	Plaintext string
	Scope     string
	Expiry    time.Time
	Hash      []byte `swaggertype:"string" format:"base64"`
	UserID    int64
}

func (t *TokenFull) Render(_ http.ResponseWriter, _ *http.Request) error {
	return nil
}

func ValidateTokenPlaintext(v *validator.Validator, tokenPlaintext string) {
	v.Check(tokenPlaintext != emptyString, keyToken, "must be provided")
	v.Check(len(tokenPlaintext) == lengthToken, keyToken, "must be 26 bytes long")
}

func (q *Queries) CreateTokenHelper(ctx context.Context, uid int64, ttl time.Duration, s string) (TokenFull, error) {
	randomBytes := make([]byte, lenthRandom)

	_, err := rand.Read(randomBytes)
	if err != nil {
		return TokenFull{}, fmt.Errorf("failed read rand: %w", err)
	}

	plaintext := base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(randomBytes)
	hash := sha256.Sum256([]byte(plaintext))

	token, err := q.CreateToken(ctx, CreateTokenParams{
		Hash:   hash[:],
		UserID: uid,
		Expiry: time.Now().Add(ttl),
		Scope:  s,
	})
	if err != nil {
		return TokenFull{}, fmt.Errorf("failed create token: %w", err)
	}

	tokenPlaintext := TokenFull{
		Plaintext: plaintext,
		Hash:      token.Hash,
		UserID:    token.UserID,
		Expiry:    token.Expiry,
		Scope:     token.Scope,
	}

	return tokenPlaintext, nil
}
