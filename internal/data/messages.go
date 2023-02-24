package data

import "github.com/seanflannery10/core/pkg/validator"

func ValidateMessage(v *validator.Validator, message string) {
	v.Check(message != "", "message", "must be provided")
	v.Check(len(message) <= 512, "message", "must not be more than 512 bytes long")
}
