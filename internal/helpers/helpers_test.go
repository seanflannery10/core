package helpers

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"

	"github.com/seanflannery10/core/internal/assert"
	"github.com/seanflannery10/core/internal/validator"
)

var ctx = context.Background()

func TestDecode(t *testing.T) {
	tests := []struct {
		a string
		e string
	}{
		{
			``,
			"body must not be empty",
		},
		{
			`<?xml version="1.0">`,
			"body contains badly-formed encode (at character 1)",
		},
		{
			`{"string": "test", }`,
			"body contains badly-formed encode (at character 20)",
		},
		{
			`["foo", "bar"]`,
			"body contains incorrect encode type (at character 1)",
		},
		{
			`{"string": 123}`,
			"body contains incorrect encode type for field \"String\"",
		},
		{
			`{"int": "123"}`,
			"body contains incorrect encode type for field \"Int\"",
		},
		{
			`{"test": 123}`,
			"body contains unknown key \"test\"",
		},
		{
			`{"int": 123}{"int": 123}`,
			"body must only contain a single encode value",
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%#v", tt.a), func(t *testing.T) {
			rr := httptest.NewRecorder()

			json := strings.NewReader(tt.a)

			r, err := http.NewRequestWithContext(ctx, http.MethodGet, "/", json)
			if err != nil {
				t.Fatal(err)
			}

			var testData struct {
				String string `a:"string"`
				Int    int    `a:"int"`
			}

			err = ReadJSON(rr, r, &testData)
			assert.Contains(t, err.Error(), tt.e)
		})
	}

	t.Run("Good", func(t *testing.T) {
		rr := httptest.NewRecorder()

		json := strings.NewReader(`{"int": 123}`)

		r, err := http.NewRequestWithContext(ctx, http.MethodGet, "/", json)
		if err != nil {
			t.Fatal(err)
		}

		var testData struct {
			String string `a:"string"`
			Int    int    `a:"int"`
		}

		err = ReadJSON(rr, r, &testData)
		if err != nil {
			t.Fatal(err)
		}
	})
}

func TestEncode(t *testing.T) {
	tests := []struct {
		c int
		s string
	}{
		{
			http.StatusOK,
			"Test 200",
		},
		{
			http.StatusNotFound,
			"Test 404",
		},
		{
			http.StatusInternalServerError,
			"Test 500",
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%#v", tt.s), func(t *testing.T) {
			rr := httptest.NewRecorder()

			err := WriteJSON(rr, tt.c, tt.s)
			if err != nil {
				t.Fatal(err)
			}

			bodyString := strings.TrimSuffix(rr.Body.String(), "\"\n")
			bodyString = strings.TrimPrefix(bodyString, "\"")

			assert.Equal(t, bodyString, tt.s)
			assert.Equal(t, rr.Header().Get("Content-Type"), "application/a")
			assert.Equal(t, rr.Code, tt.c)
		})
	}
}

func TestEncodeWithHeaders(t *testing.T) {
	tests := []struct {
		c int
		s string
		h string
	}{
		{
			http.StatusOK,
			"Test 200",
			"Test-Header",
		},
		{
			http.StatusNotFound,
			"Test 404",
			"123",
		},
		{
			http.StatusInternalServerError,
			"Test 500",
			"",
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%#v", tt.s), func(t *testing.T) {
			rr := httptest.NewRecorder()

			headers := make(http.Header)
			headers.Set("X-Request-Id", tt.h)

			err := WriteJSONWithHeaders(rr, tt.c, tt.s, headers)
			if err != nil {
				t.Fatal(err)
			}

			bodyString := strings.TrimSuffix(rr.Body.String(), "\"\n")
			bodyString = strings.TrimPrefix(bodyString, "\"")

			assert.Equal(t, bodyString, tt.s)
			assert.Equal(t, rr.Header().Get("Content-Type"), "application/a")
			assert.Equal(t, rr.Header().Get("X-Request-Id"), tt.h)
			assert.Equal(t, rr.Code, tt.c)
		})
	}
}

func TestCSV(t *testing.T) {
	tests := []struct {
		key string
		csv string
	}{
		{
			"csv",
			"csv,foo,bar",
		},
		{
			"csv",
			"csv,foo,bar,test",
		},
		{
			"string",
			"string",
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%#v", tt.csv), func(t *testing.T) {
			r, err := http.NewRequestWithContext(ctx, http.MethodGet, "/", nil)
			if err != nil {
				t.Fatal(err)
			}

			split := strings.Split(tt.csv, ",")

			r.URL.RawQuery = url.Values{
				tt.key: {tt.csv},
			}.Encode()

			qs := r.URL.Query()
			res := ReadCSVParam(qs, tt.key, nil)

			assert.Equal(t, len(res), len(split))
			assert.Equal(t, res[0], split[0])
			assert.Equal(t, res[len(res)-1], split[len(split)-1])
		})
	}
}

func TestInt(t *testing.T) {
	tests := []struct {
		key string
		int int
	}{
		{
			"test",
			42,
		},
		{
			"test",
			0,
		},
		{
			"test",
			-20,
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%#v", tt.int), func(t *testing.T) {
			r, err := http.NewRequestWithContext(ctx, http.MethodGet, "/", nil)
			if err != nil {
				t.Fatal(err)
			}

			r.URL.RawQuery = url.Values{
				tt.key: {strconv.Itoa(tt.int)},
			}.Encode()

			qs := r.URL.Query()
			res := ReadIntParam(qs, tt.key, 1, &validator.Validator{})

			assert.Equal(t, res, tt.int)
		})
	}
}

func TestString(t *testing.T) {
	tests := []struct {
		key    string
		string string
	}{
		{
			"test",
			"42",
		},
		{
			"test",
			"test",
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%#v", tt.string), func(t *testing.T) {
			r, err := http.NewRequestWithContext(ctx, http.MethodGet, "/", nil)
			if err != nil {
				t.Fatal(err)
			}

			r.URL.RawQuery = url.Values{
				tt.key: {tt.string},
			}.Encode()

			qs := r.URL.Query()
			res := ReadStringParam(qs, tt.key, "")

			assert.Equal(t, res, tt.string)
		})
	}
}
