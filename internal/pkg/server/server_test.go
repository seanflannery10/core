package server_test

import (
	"os"
	"syscall"
	"testing"
	"time"

	"github.com/seanflannery10/core/internal/pkg/server"
)

func TestServer_Run(t *testing.T) {
	t.Run("SIGINT", func(t *testing.T) {
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

		err := server.Serve(4444, nil)
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("SIGTERM", func(t *testing.T) {
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

		err := server.Serve(4444, nil)
		if err != nil {
			t.Fatal(err)
		}
	})
}
