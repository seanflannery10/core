// Code generated by ogen, DO NOT EDIT.

package api

import (
	"github.com/go-faster/errors"

	"github.com/ogen-go/ogen/validate"
)

func (s *MessageRequest) Validate() error {
	var failures []validate.FieldError
	if err := func() error {
		if err := (validate.String{
			MinLength:    0,
			MinLengthSet: false,
			MaxLength:    512,
			MaxLengthSet: true,
			Email:        false,
			Hostname:     false,
			Regex:        nil,
		}).Validate(string(s.Message)); err != nil {
			return errors.Wrap(err, "string")
		}
		return nil
	}(); err != nil {
		failures = append(failures, validate.FieldError{
			Name:  "message",
			Error: err,
		})
	}
	if len(failures) > 0 {
		return &validate.Error{Fields: failures}
	}
	return nil
}
func (s *MessagesResponse) Validate() error {
	var failures []validate.FieldError
	if err := func() error {
		if s.Messages == nil {
			return errors.New("nil is invalid value")
		}
		return nil
	}(); err != nil {
		failures = append(failures, validate.FieldError{
			Name:  "messages",
			Error: err,
		})
	}
	if len(failures) > 0 {
		return &validate.Error{Fields: failures}
	}
	return nil
}

func (s *TokenRequest) Validate() error {
	var failures []validate.FieldError
	if err := func() error {
		if err := (validate.String{
			MinLength:    26,
			MinLengthSet: true,
			MaxLength:    26,
			MaxLengthSet: true,
			Email:        false,
			Hostname:     false,
			Regex:        nil,
		}).Validate(string(s.Token)); err != nil {
			return errors.Wrap(err, "string")
		}
		return nil
	}(); err != nil {
		failures = append(failures, validate.FieldError{
			Name:  "token",
			Error: err,
		})
	}
	if len(failures) > 0 {
		return &validate.Error{Fields: failures}
	}
	return nil
}
func (s *UpdateUserPasswordRequest) Validate() error {
	var failures []validate.FieldError
	if err := func() error {
		if err := (validate.String{
			MinLength:    8,
			MinLengthSet: true,
			MaxLength:    72,
			MaxLengthSet: true,
			Email:        false,
			Hostname:     false,
			Regex:        nil,
		}).Validate(string(s.Password)); err != nil {
			return errors.Wrap(err, "string")
		}
		return nil
	}(); err != nil {
		failures = append(failures, validate.FieldError{
			Name:  "password",
			Error: err,
		})
	}
	if err := func() error {
		if err := (validate.String{
			MinLength:    26,
			MinLengthSet: true,
			MaxLength:    26,
			MaxLengthSet: true,
			Email:        false,
			Hostname:     false,
			Regex:        nil,
		}).Validate(string(s.Token)); err != nil {
			return errors.Wrap(err, "string")
		}
		return nil
	}(); err != nil {
		failures = append(failures, validate.FieldError{
			Name:  "token",
			Error: err,
		})
	}
	if len(failures) > 0 {
		return &validate.Error{Fields: failures}
	}
	return nil
}
func (s *UserEmailRequest) Validate() error {
	var failures []validate.FieldError
	if err := func() error {
		if err := (validate.String{
			MinLength:    0,
			MinLengthSet: false,
			MaxLength:    0,
			MaxLengthSet: false,
			Email:        true,
			Hostname:     false,
			Regex:        nil,
		}).Validate(string(s.Email)); err != nil {
			return errors.Wrap(err, "string")
		}
		return nil
	}(); err != nil {
		failures = append(failures, validate.FieldError{
			Name:  "email",
			Error: err,
		})
	}
	if len(failures) > 0 {
		return &validate.Error{Fields: failures}
	}
	return nil
}
func (s *UserLoginRequest) Validate() error {
	var failures []validate.FieldError
	if err := func() error {
		if err := (validate.String{
			MinLength:    0,
			MinLengthSet: false,
			MaxLength:    0,
			MaxLengthSet: false,
			Email:        true,
			Hostname:     false,
			Regex:        nil,
		}).Validate(string(s.Email)); err != nil {
			return errors.Wrap(err, "string")
		}
		return nil
	}(); err != nil {
		failures = append(failures, validate.FieldError{
			Name:  "email",
			Error: err,
		})
	}
	if err := func() error {
		if err := (validate.String{
			MinLength:    8,
			MinLengthSet: true,
			MaxLength:    72,
			MaxLengthSet: true,
			Email:        false,
			Hostname:     false,
			Regex:        nil,
		}).Validate(string(s.Password)); err != nil {
			return errors.Wrap(err, "string")
		}
		return nil
	}(); err != nil {
		failures = append(failures, validate.FieldError{
			Name:  "password",
			Error: err,
		})
	}
	if len(failures) > 0 {
		return &validate.Error{Fields: failures}
	}
	return nil
}
func (s *UserRequest) Validate() error {
	var failures []validate.FieldError
	if err := func() error {
		if err := (validate.String{
			MinLength:    0,
			MinLengthSet: false,
			MaxLength:    100,
			MaxLengthSet: true,
			Email:        false,
			Hostname:     false,
			Regex:        nil,
		}).Validate(string(s.Name)); err != nil {
			return errors.Wrap(err, "string")
		}
		return nil
	}(); err != nil {
		failures = append(failures, validate.FieldError{
			Name:  "name",
			Error: err,
		})
	}
	if err := func() error {
		if err := (validate.String{
			MinLength:    0,
			MinLengthSet: false,
			MaxLength:    0,
			MaxLengthSet: false,
			Email:        true,
			Hostname:     false,
			Regex:        nil,
		}).Validate(string(s.Email)); err != nil {
			return errors.Wrap(err, "string")
		}
		return nil
	}(); err != nil {
		failures = append(failures, validate.FieldError{
			Name:  "email",
			Error: err,
		})
	}
	if err := func() error {
		if err := (validate.String{
			MinLength:    8,
			MinLengthSet: true,
			MaxLength:    72,
			MaxLengthSet: true,
			Email:        false,
			Hostname:     false,
			Regex:        nil,
		}).Validate(string(s.Password)); err != nil {
			return errors.Wrap(err, "string")
		}
		return nil
	}(); err != nil {
		failures = append(failures, validate.FieldError{
			Name:  "password",
			Error: err,
		})
	}
	if len(failures) > 0 {
		return &validate.Error{Fields: failures}
	}
	return nil
}
func (s *UserResponse) Validate() error {
	var failures []validate.FieldError
	if err := func() error {
		if err := (validate.String{
			MinLength:    0,
			MinLengthSet: false,
			MaxLength:    0,
			MaxLengthSet: false,
			Email:        true,
			Hostname:     false,
			Regex:        nil,
		}).Validate(string(s.Email)); err != nil {
			return errors.Wrap(err, "string")
		}
		return nil
	}(); err != nil {
		failures = append(failures, validate.FieldError{
			Name:  "email",
			Error: err,
		})
	}
	if len(failures) > 0 {
		return &validate.Error{Fields: failures}
	}
	return nil
}
