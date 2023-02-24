package data

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/seanflannery10/core/pkg/validator"
)

const (
	ScopeActivation     = "activation"
	ScopeAuthentication = "authentication"
	ScopePasswordReset  = "password-reset"
)

type FullToken struct {
	Plaintext string
	Hash      []byte
	UserID    int64
	Expiry    pgtype.Timestamptz
	Scope     string
}

func (q *Queries) NewToken(ctx context.Context, userID int64, ttl time.Duration, scope string) (FullToken, error) {
	randomBytes := make([]byte, 16)

	_, err := rand.Read(randomBytes)
	if err != nil {
		return FullToken{}, err
	}

	plaintext := base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(randomBytes)
	hash := sha256.Sum256([]byte(plaintext))

	token, err := q.CreateToken(ctx, CreateTokenParams{
		Hash:   hash[:],
		UserID: userID,
		Expiry: pgtype.Timestamptz{Time: time.Now().Add(ttl), Valid: true},
		Scope:  scope,
	})
	if err != nil {
		return FullToken{}, err
	}

	fullToken := FullToken{
		Plaintext: plaintext,
		Hash:      token.Hash,
		UserID:    token.UserID,
		Expiry:    token.Expiry,
		Scope:     token.Scope,
	}

	return fullToken, nil
}

func ValidateTokenPlaintext(v *validator.Validator, tokenPlaintext string) {
	v.Check(tokenPlaintext != "", "token", "must be provided")
	v.Check(len(tokenPlaintext) == 26, "token", "must be 26 bytes long")
}
