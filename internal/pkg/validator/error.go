package validator

var ErrValidation *ValidationError

type ValidationError struct {
	Errors map[string]string
}

func NewValidationError(errors map[string]string) *ValidationError {
	return &ValidationError{
		Errors: errors,
	}
}

func (ve *ValidationError) Error() string {
	return "validation error"
}
