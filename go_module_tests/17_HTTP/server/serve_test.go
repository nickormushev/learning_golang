package server_test

import (
	server "learning/17_HTTP/server"
	"syscall"
	"testing"
	"time"
)

func TestCreateContext(t *testing.T) {
	t.Run("Context should be killed only after SIGINT", func(t *testing.T) {
		ctx := server.GenerateContextWithSigint()

		select {
		case <-ctx.Done():
			t.Errorf("Expected context to not be cancelled before SIGINT but it was")
		case <-time.After(5 * time.Millisecond):
		}

		syscall.Kill(syscall.Getpid(), syscall.SIGINT)

		select {
		case <-ctx.Done():
		case <-time.After(5 * time.Millisecond):
			t.Errorf("Expected context be cancelled afte SIGIN but was not")
		}
	})
}
