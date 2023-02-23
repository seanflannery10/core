package validator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidator(t *testing.T) {
	v := New()

	assert.Equal(t, v.HasErrors(), false)

	v.AddError("test", "test field error")
	assert.Equal(t, v.Errors["test"], "test field error")

	v.Check(true, "test2", "test field error2")
	assert.Equal(t, len(v.Errors), 1)

	v.Check(false, "test3", "test field error3")
	assert.Equal(t, v.Errors["test3"], "test field error3")

	assert.Equal(t, v.HasErrors(), true)
}
