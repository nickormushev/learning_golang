package main

import (
	"fmt"
	"io"
	"os"
	"time"
)

//Sleeper is an interface for struct that have a Sleep() method
type Sleeper interface {
	Sleep()
}

//ConfigurableSleeper is a Sleeper that can be modified to sleep using different functions for different durations
type ConfigurableSleeper struct {
	sleepTime time.Duration
	sleepFunc func(time.Duration)
}

//Sleep calls a sleeperfunction with a given time.Duration
func (c ConfigurableSleeper) Sleep() {
	c.sleepFunc(c.sleepTime)
}

//DefaultSleeper is used to simulate an object we are mocking. It sleeps for one second
type DefaultSleeper struct{}

//Sleep sleeps for one second
func (d DefaultSleeper) Sleep() {
	time.Sleep(1 * time.Second)
}

const countdownStart int = 3

//Countdown counts from 3 to 1 and writes Go! afterwards
func Countdown(w io.Writer, sleeper Sleeper) {
	sleeper.Sleep()
	for i := countdownStart; i > 0; i-- {
		fmt.Fprintln(w, i)
		sleeper.Sleep()
	}

	fmt.Fprint(w, "Go!")
}

func main() {
	confSleeper := ConfigurableSleeper{1 * time.Second, time.Sleep}
	Countdown(os.Stdout, confSleeper)
}
