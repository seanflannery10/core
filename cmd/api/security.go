package main

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"strings"
	"time"

	"github.com/go-faster/errors"
	"github.com/seanflannery10/core/internal/api"
	"github.com/seanflannery10/core/internal/data"
	"github.com/seanflannery10/core/internal/shared/utils"
	"github.com/segmentio/asm/base64"
)

var (
	errInvalidAccessToken = errors.New("invalid access token")
	errUserNotActivated   = errors.New("user account must be activated")
)

type security struct {
	Queries   *data.Queries
	SecretKey []byte
}

func (s *security) HandleAccess(ctx context.Context, _ string, t api.Access) (context.Context, error) {
	tokenHash := sha256.Sum256([]byte(t.Token))

	user, err := s.Queries.GetUserFromToken(ctx, data.GetUserFromTokenParams{
		Hash:   tokenHash[:],
		Scope:  data.ScopeAccess,
		Expiry: time.Now(),
	})
	if err != nil {
		return ctx, errInvalidAccessToken
	}

	if !user.Activated {
		return ctx, errUserNotActivated
	}

	return utils.ContextSetUser(ctx, &user), nil
}

func (s *security) HandleRefresh(ctx context.Context, _ string, r api.Refresh) (context.Context, error) {
	encryptedValue, err := base64.URLEncoding.DecodeString(r.APIKey)
	if err != nil {
		return ctx, errors.Wrap(err, "failed decode string")
	}

	block, err := aes.NewCipher(s.SecretKey)
	if err != nil {
		return ctx, errors.Wrap(err, "failed cipher")
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return ctx, errors.Wrap(err, "failed new gcm")
	}

	nonceSize := aesGCM.NonceSize()

	if len(encryptedValue) < nonceSize {
		return ctx, nil
	}

	nonce := encryptedValue[:nonceSize]
	ciphertext := encryptedValue[nonceSize:]

	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return ctx, errors.Wrap(err, "failed gcm open")
	}

	_, value, ok := strings.Cut(string(plaintext), ":")
	if !ok {
		return ctx, errors.Wrap(err, "failed cut string")
	}

	return utils.ContextSetCookieValue(ctx, value), nil
}
