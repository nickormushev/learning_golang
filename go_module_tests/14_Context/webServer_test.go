package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/pkg/errors"
)

type SpyStore struct {
	response  string
	cancelled bool
	t         *testing.T
}

func (s *SpyStore) Cancel() {
	s.cancelled = true
}

func (s *SpyStore) Fetch(ctx context.Context) (string, error) {
	ch := make(chan string, 1)

	go func() {
		var result string
		for _, c := range s.response {
			select {
			case <-ctx.Done():
				s.t.Log("Fetch cancelled")
				return
			default:
				time.Sleep(time.Millisecond * 10)
				result += string(c)
			}
		}

		ch <- result
	}()

	select {
	case data := <-ch:
		return data, nil
	case <-ctx.Done():
		return "", ctx.Err()
	}
}

func (s *SpyStore) assertWasNotCancelled() {
	s.t.Helper()
	if s.cancelled {
		s.t.Errorf("Request was cancelled but shouldn't have been")
	}
}

func (s *SpyStore) assertWasCancelled() {
	s.t.Helper()
	if !s.cancelled {
		s.t.Errorf("Request was cancelled but shouldn't have been")
	}
}

type SpyResponseWriter struct {
	written bool
}

func (s *SpyResponseWriter) Header() http.Header {
	s.written = true
	return nil
}

func (s *SpyResponseWriter) Write([]byte) (int, error) {
	s.written = true
	return 0, errors.New("Not implemented")
}

func (s *SpyResponseWriter) WriteHeader(statusCode int) {
	s.written = true
}

func TestHandler(t *testing.T) {
	data := "Hello, Gosho"
	t.Run("Successfull response", func(t *testing.T) {
		spyStore := &SpyStore{response: data, t: t}
		srv := Server(spyStore)

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		response := httptest.NewRecorder()

		srv.ServeHTTP(response, req)

		//spyStore.assertWasNotCancelled()

		if response.Body.String() != data {
			t.Errorf("Expected response: %s but got %s", data, response.Body.String())
		}
	})

	t.Run("Cancel due to long Fetch()", func(t *testing.T) {
		spyStore := &SpyStore{response: data, t: t}
		srv := Server(spyStore)

		req := httptest.NewRequest(http.MethodGet, "/", nil)

		cancellingCtx, cancel := context.WithCancel(req.Context())
		time.AfterFunc(5*time.Millisecond, cancel)
		req = req.WithContext(cancellingCtx)

		response := &SpyResponseWriter{}

		srv.ServeHTTP(response, req)

		if response.written {
			t.Errorf("A response should not have been written")
		}
	})
}
