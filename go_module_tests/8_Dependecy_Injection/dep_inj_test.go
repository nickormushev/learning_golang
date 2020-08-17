package main

import (
	"bytes"
	"testing"
)

func TestGreet(t *testing.T) {
	buff := bytes.Buffer{}

	Greet(&buff, "Joe")
	got := buff.String()
	want := "Hello, Joe!\n"

	if got != want {
		t.Errorf("wanted: %s, got: %s", want, got)
	}
}
