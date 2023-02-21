package server

import (
	"os"
	"syscall"
	"testing"
	"time"

	"github.com/seanflannery10/core/internal/assert"
)

func TestNew(t *testing.T) {
	srv := New(12345, nil)
	srv.Background(func() {})

	assert.Equal(t, srv.Addr, ":12345")
	assert.SameType(t, srv, &Server{})
}

func TestServer_Run(t *testing.T) {
	t.Run("SIGINT", func(t *testing.T) {
		srv := New(4444, nil)

		go func() {
			time.Sleep(250 * time.Millisecond)

			p, err := os.FindProcess(os.Getpid())
			if err != nil {
				panic(err)
			}

			err = p.Signal(syscall.SIGINT)
			if err != nil {
				return
			}
		}()

		err := srv.Serve()
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("SIGTERM", func(t *testing.T) {
		srv := New(4444, nil)

		go func() {
			time.Sleep(250 * time.Millisecond)

			p, err := os.FindProcess(os.Getpid())
			if err != nil {
				panic(err)
			}

			err = p.Signal(syscall.SIGTERM)
			if err != nil {
				return
			}
		}()

		err := srv.Serve()
		if err != nil {
			t.Fatal(err)
		}
	})
}
