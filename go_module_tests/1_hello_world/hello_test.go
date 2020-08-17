package main

import "testing"

func TestHello(t *testing.T) {
	assertCorrectMessage := func(t *testing.T, got, want string) {
		t.Helper()
		if got != want {
			t.Errorf("Failed to say %q and said %q", want, got)
		}
	}

	t.Run("Saying hello to the world", func(t *testing.T) {
		got := Hello("")
		want := "Hello, world!\n"

		assertCorrectMessage(t, got, want)
	})

	t.Run("Saying hello to Pesho", func(t *testing.T) {
		got := Hello("Pesho")
		want := "Hello, Pesho!\n"

		assertCorrectMessage(t, got, want)
	})
}
