package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestRacer(t *testing.T) {
	t.Run("Returns the faster server", func(t *testing.T) {
		slowServer := makeDelayedServer(10 * time.Millisecond)
		fastServer := makeDelayedServer(5 * time.Millisecond)

		defer fastServer.Close()
		defer slowServer.Close()

		slowURL := slowServer.URL
		fastURL := fastServer.URL

		want := fastURL
		got, err := Racer(slowURL, fastURL)

		if err != nil {
			t.Fatal("Got and error but expected nil")
		}

		if got != want {
			t.Errorf("wanted %s but got %s", want, got)
		}
	})

	t.Run("Returns and error when waiting more than 10 seconds", func(t *testing.T) {
		serverA := makeDelayedServer(11 * time.Millisecond)
		serverB := makeDelayedServer(11 * time.Millisecond)

		defer serverA.Close()
		defer serverB.Close()

		_, err := ConfigurableRacer(serverA.URL, serverB.URL, 10*time.Millisecond)

		if err == nil {
			t.Fatal("Expected errro but got nil")
		}
	})
}

func makeDelayedServer(t time.Duration) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(t)
		w.WriteHeader(http.StatusOK)
	}))
}
