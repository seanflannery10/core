package cookies

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

var (
	ErrValueTooLong = errors.New("cookie value too long")
	ErrInvalidValue = errors.New("invalid cookie value")
)

const (
	cookieMaxSize = 4096
	emptyString   = ""
)

func WriteEncrypted(w http.ResponseWriter, cookie *http.Cookie, secretKey []byte) error {
	block, err := aes.NewCipher(secretKey)
	if err != nil {
		return fmt.Errorf("failed new cipher: %w", err)
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return fmt.Errorf("failed new gcm: %w", err)
	}

	nonce := make([]byte, aesGCM.NonceSize())

	_, err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return fmt.Errorf("failed read full: %w", err)
	}

	plaintext := fmt.Sprintf("%s:%s", cookie.Name, cookie.Value)

	encryptedValue := aesGCM.Seal(nonce, nonce, []byte(plaintext), nil)

	cookie.Value = base64.URLEncoding.EncodeToString(encryptedValue)

	if len(cookie.String()) > cookieMaxSize {
		return fmt.Errorf("failed length check: %w", ErrValueTooLong)
	}

	http.SetCookie(w, cookie)

	return nil
}

func ReadEncrypted(r *http.Request, name string, secretKey []byte) (string, error) {
	cookie, err := r.Cookie(name)
	if err != nil {
		return emptyString, fmt.Errorf("failed read cookie: %w", err)
	}

	encryptedValue, err := base64.URLEncoding.DecodeString(cookie.Value)
	if err != nil {
		return emptyString, fmt.Errorf("failed decode string: %w", ErrInvalidValue)
	}

	plaintext, err := getPlainText(encryptedValue, secretKey)
	if err != nil {
		return emptyString, fmt.Errorf("failed get plain text: %w", ErrInvalidValue)
	}

	expectedName, value, ok := strings.Cut(string(plaintext), ":")
	if !ok {
		return emptyString, fmt.Errorf("failed string cut: %w", ErrInvalidValue)
	}

	if expectedName != name {
		return emptyString, fmt.Errorf("failed name check: %w", ErrInvalidValue)
	}

	return value, nil
}

func getPlainText(encryptedValue, secretKey []byte) ([]byte, error) {
	block, err := aes.NewCipher(secretKey)
	if err != nil {
		return nil, fmt.Errorf("failed new cipher: %w", err)
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed new gcm: %w", err)
	}

	nonceSize := aesGCM.NonceSize()

	if len(encryptedValue) < nonceSize {
		return nil, fmt.Errorf("failed nonce size: %w", ErrInvalidValue)
	}

	nonce := encryptedValue[:nonceSize]
	ciphertext := encryptedValue[nonceSize:]

	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("failed gcm open: %w", ErrInvalidValue)
	}

	return plaintext, nil
}
