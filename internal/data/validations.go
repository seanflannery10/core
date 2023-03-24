package data

import "github.com/seanflannery10/core/internal/shared/validator"

const (
	emptyString   = ""
	keyToken      = "token"
	lengthToken   = 26
	keyEmail      = "email"
	keyPassword   = "password"
	keyName       = "name"
	requiredField = "must be provided"
)

func ValidateMessage(v *validator.Validator, message string) {
	v.Check(message != "", "message", "must be provided")
	v.Check(len(message) <= 512, "message", "must not be more than 512 bytes long")
}

func ValidateTokenPlaintext(v *validator.Validator, tokenPlaintext string) {
	v.Check(tokenPlaintext != emptyString, keyToken, "must be provided")
	v.Check(len(tokenPlaintext) == lengthToken, keyToken, "must be 26 bytes long")
}

func ValidateEmail(v *validator.Validator, email string) {
	v.Check(email != "", keyEmail, requiredField)
	v.Check(validator.Matches(email, validator.RgxEmail), keyEmail, "must be a valid email address")
}

func ValidatePasswordPlaintext(v *validator.Validator, password string) {
	v.Check(password != "", keyPassword, requiredField)
	v.Check(len(password) >= 8, keyPassword, "must be at least 8 bytes long")
	v.Check(len(password) <= 72, keyPassword, "must not be more than 72 bytes long")
}

func ValidateName(v *validator.Validator, name string) {
	v.Check(name != "", keyName, requiredField)
	v.Check(len(name) <= 500, keyName, "must not be more than 500 bytes long")
}
