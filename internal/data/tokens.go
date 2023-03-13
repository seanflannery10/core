package data

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/seanflannery10/core/internal/pkg/validator"
)

const (
	ScopeActivation    = "activation"
	ScopeAccess        = "access"
	ScopePasswordReset = "password-reset"
	ScopeRefresh       = "refresh"
)

type TokenFull struct {
	Plaintext string
	Hash      []byte
	UserID    int64
	Expiry    pgtype.Timestamptz
	Scope     string
}

func (t *TokenFull) Render(_ http.ResponseWriter, _ *http.Request) error {
	return nil
}

func ValidateTokenPlaintext(v *validator.Validator, tokenPlaintext string) {
	v.Check(tokenPlaintext != "", "token", "must be provided")
	v.Check(len(tokenPlaintext) == 26, "token", "must be 26 bytes long")
}

func (q *Queries) CreateTokenHelper(ctx context.Context, uid int64, ttl time.Duration, s string) (TokenFull, error) {
	randomBytes := make([]byte, 16)

	_, err := rand.Read(randomBytes)
	if err != nil {
		return TokenFull{}, err
	}

	plaintext := base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(randomBytes)
	hash := sha256.Sum256([]byte(plaintext))

	token, err := q.CreateToken(ctx, CreateTokenParams{
		Hash:   hash[:],
		UserID: uid,
		Expiry: pgtype.Timestamptz{Time: time.Now().Add(ttl), Valid: true},
		Scope:  s,
	})
	if err != nil {
		return TokenFull{}, err
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
