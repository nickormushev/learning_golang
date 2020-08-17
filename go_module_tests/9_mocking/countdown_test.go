package main

import (
	"bytes"
	"reflect"
	"testing"
	"time"
)

const (
	sleep = "sleep"
	write = "write"
)

type SpySleeper struct {
	Calls int
}

type SpyTime struct {
	durationSlept time.Duration
}

func (s *SpyTime) Sleep(duration time.Duration) {
	s.durationSlept += duration
}

type CountdownOperationsSpy struct {
	Calls []string
}

func (c *CountdownOperationsSpy) Sleep() {
	c.Calls = append(c.Calls, sleep)
}

func (c *CountdownOperationsSpy) Write(p []byte) (n int, err error) {
	c.Calls = append(c.Calls, write)
	return
}

func (s *SpySleeper) Sleep() {
	s.Calls++
}

func TestConfigurableSleeper(t *testing.T) {
	buff := &bytes.Buffer{}
	spySleeper := &SpyTime{}
	wantSleep := 6 * time.Second
	confSleeper := ConfigurableSleeper{time.Duration(wantSleep), spySleeper.Sleep}
	Countdown(buff, confSleeper)

	if spySleeper.durationSlept != 4*wantSleep {
		t.Errorf("Wanted to sleep %d but slept %d", wantSleep, spySleeper.durationSlept)
	}
}

func TestCountdown(t *testing.T) {
	t.Run("Sleep count and output test", func(t *testing.T) {
		buff := &bytes.Buffer{}
		spySleeper := &SpySleeper{}

		Countdown(buff, spySleeper)

		got := buff.String()
		want := `3
2
1
Go!`

		if got != want {
			t.Errorf("Wanted %s but got %s", want, got)
		}

		if spySleeper.Calls != 4 {
			t.Errorf("Sleep function is not called 3 times but %d!", spySleeper.Calls)
		}
	})

	t.Run("sleep before every print", func(t *testing.T) {
		spySleepPrinter := &CountdownOperationsSpy{}
		Countdown(spySleepPrinter, spySleepPrinter)

		want := []string{
			sleep,
			write,
			sleep,
			write,
			sleep,
			write,
			sleep,
			write,
		}

		if !reflect.DeepEqual(want, spySleepPrinter.Calls) {
			t.Errorf("Wanted calls %v times but got %v!", want, spySleepPrinter.Calls)
		}
	})
}
