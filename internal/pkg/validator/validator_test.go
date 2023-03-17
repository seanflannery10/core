package validator_test

import (
	"testing"

	"github.com/seanflannery10/core/internal/pkg/assert"
	"github.com/seanflannery10/core/internal/pkg/validator"
)

const (
	equalValue = 1
)

func TestValidator(t *testing.T) {
	v := validator.New()

	assert.Equal(t, v.HasErrors(), false)

	v.AddError("test", "test field error")
	assert.Equal(t, v.Errors["test"], "test field error")

	v.Check(true, "test2", "test field error2")
	assert.Equal(t, len(v.Errors), equalValue)

	v.Check(false, "test3", "test field error3")
	assert.Equal(t, v.Errors["test3"], "test field error3")

	assert.Equal(t, v.HasErrors(), true)
}
