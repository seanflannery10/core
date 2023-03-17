package assert_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/seanflannery10/core/internal/pkg/assert"
)

const (
	d2000       = 200.0
	d2001       = 200.1
	emptyString = ""
	i0          = 0
	i200        = 200
	i201        = 201
	testString  = "test"
	v           = "%#v"
	string1     = "123"
	string2     = "124"
)

//nolint:all
func TestEqual(t *testing.T) {
	testsNums := []struct {
		a float32
		e float32
		r bool
	}{
		{
			i0,
			i0,
			true,
		},
		{
			i200,
			i200,
			true,
		},
		{
			i201,
			i200,
			false,
		},
		{
			d2000,
			d2000,
			true,
		},
		{
			d2001,
			d2000,
			false,
		},
		{
			-i200,
			-i200,
			true,
		},
		{
			-i201,
			-i200,
			false,
		},
		{
			-d2000,
			-d2000,
			true,
		},
		{
			-d2001,
			-d2000,
			false,
		},
	}

	for _, tt := range testsNums {
		t.Run(fmt.Sprintf(v, tt.a), func(t *testing.T) {
			res := assert.Equal(new(testing.T), tt.a, tt.e)
			if res != tt.r {
				t.Errorf("Equal(%#v, %#v) should return %#v", tt.a, tt.e, tt.r)
			}
		})
	}

	testsStrings := []struct {
		a string
		e string
		r bool
	}{
		{
			string1,
			string1,
			true,
		},
		{
			string1,
			string2,
			false,
		},
		{
			emptyString,
			emptyString,
			true,
		},
	}

	for _, tt := range testsStrings {
		t.Run(fmt.Sprintf(v, tt.a), func(t *testing.T) {
			res := assert.Equal(new(testing.T), tt.a, tt.e)
			if res != tt.r {
				t.Errorf("Equal(%#v, %#v) should return %#v", tt.a, tt.e, tt.r)
			}
		})
	}
}

//nolint:all
func TestNotEqual(t *testing.T) {
	testsNums := []struct {
		a float32
		e float32
		r bool
	}{
		{
			i0,
			i0,
			false,
		},
		{
			i200,
			i200,
			false,
		},
		{
			i201,
			i200,
			true,
		},
		{
			d2000,
			d2000,
			false,
		},
		{
			d2001,
			d2000,
			true,
		},
		{
			-i200,
			-i200,
			false,
		},
		{
			-i201,
			-i200,
			true,
		},
		{
			-d2000,
			-d2000,
			false,
		},
		{
			-d2001,
			-d2000,
			true,
		},
	}

	for _, tt := range testsNums {
		t.Run(fmt.Sprintf(v, tt.a), func(t *testing.T) {
			res := assert.NotEqual(new(testing.T), tt.a, tt.e)
			if res != tt.r {
				t.Errorf("NotEqual(%#v, %#v) should return %#v", tt.a, tt.e, tt.r)
			}
		})
	}

	testsStrings := []struct {
		a string
		e string
		r bool
	}{
		{
			string1,
			string1,
			false,
		},
		{
			string2,
			string1,
			true,
		},
		{
			emptyString,
			emptyString,
			false,
		},
	}

	for _, tt := range testsStrings {
		t.Run(fmt.Sprintf(v, tt.a), func(t *testing.T) {
			res := assert.NotEqual(new(testing.T), tt.a, tt.e)
			if res != tt.r {
				t.Errorf("NotEqual(%#v, %#v) should return %#v", tt.a, tt.e, tt.r)
			}
		})
	}
}

func TestSameType(t *testing.T) {
	t.Run("String", func(t *testing.T) {
		res := assert.SameType(new(testing.T), "foo", "bar")
		if !res {
			t.Errorf("SameType(%#v, %#v) should return %#v", "foo", "bar", "true")
		}
	})

	t.Run("Int", func(t *testing.T) {
		res := assert.SameType(new(testing.T), 1, 2)
		if !res {
			t.Errorf("SameType(%#v, %#v) should return %#v", "foo", "bar", "true")
		}
	})

	t.Run("Float", func(t *testing.T) {
		res := assert.SameType(new(testing.T), 1.2, 2.3)
		if !res {
			t.Errorf("SameType(%#v, %#v) should return %#v", "foo", "bar", "true")
		}
	})

	t.Run("Bool", func(t *testing.T) {
		res := assert.SameType(new(testing.T), true, false)
		if !res {
			t.Errorf("SameType(%#v, %#v) should return %#v", "foo", "bar", "true")
		}
	})
}

func TestContains(t *testing.T) {
	tests := []struct {
		a string
		e string
		r bool
	}{
		{
			"this is a test",
			"test",
			true,
		},
		{
			"this is a",
			"test",
			false,
		},
		{
			emptyString,
			emptyString,
			true,
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf(v, tt.e), func(t *testing.T) {
			res := assert.Contains(new(testing.T), tt.a, tt.e)
			if res != tt.r {
				t.Errorf("Contains(%#v, %#v) should return %#v", tt.a, tt.e, tt.r)
			}
		})
	}
}

func TestNotContains(t *testing.T) {
	tests := []struct {
		a string
		e string
		r bool
	}{
		{
			"this is a test",
			testString,
			false,
		},
		{
			"this is a",
			testString,
			true,
		},
		{
			emptyString,
			emptyString,
			false,
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf(v, tt.e), func(t *testing.T) {
			res := assert.NotContains(new(testing.T), tt.a, tt.e)
			if res != tt.r {
				t.Errorf("NotContains(%#v, %#v) should return %#v", tt.a, tt.e, tt.r)
			}
		})
	}
}

func TestNilError(t *testing.T) {
	tests := []struct {
		a error
		r bool
	}{
		{
			nil,
			true,
		},
		{
			errors.New(testString), //nolint:goerr113
			false,
		},
		{
			errors.New(emptyString), //nolint:goerr113
			false,
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf(v, tt.a), func(t *testing.T) {
			res := assert.NilError(new(testing.T), tt.a)
			if res != tt.r {
				t.Errorf("NilError(%#v) should return %#v", tt.a, tt.r)
			}
		})
	}
}
