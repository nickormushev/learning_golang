package main

import (
	"net/http"
	"time"
)

var tenSecondTimeout time.Duration = 10 * time.Second

func getTime(url string) time.Duration {
	start := time.Now()
	http.Get(url)
	return time.Since(start)
}

//Error is a type for Errors used so const works
type Error string

func (e Error) Error() string {
	return string(e)
}

//ErrNoResponse No response from servers in Racer for more than 10 seconds
const ErrNoResponse Error = "No response after waiting for ten seconds"

//Racer issues a get request to two urls and returns the one that responded quicker
func Racer(a, b string) (fastest string, err error) {
	return ConfigurableRacer(a, b, tenSecondTimeout)
}

//ConfigurableRacer issues a get request to two urls and returns the one that responded quicker
func ConfigurableRacer(a, b string, timeout time.Duration) (fastest string, err error) {
	select {
	case <-ping(a):
		return a, nil
	case <-ping(b):
		return b, nil
	case <-time.After(timeout):
		return "", ErrNoResponse
	}

}

func ping(url string) chan struct{} {
	ch := make(chan struct{})

	go func() {
		http.Get(url)
		close(ch)
	}()

	return ch
}

//Initial implementation for reference
//func racer(a, b string) (fastest string) {
//	timea := gettime(a)
//	timeb := gettime(b)
//
//	if timea < timeb {
//		return a
//	}
//
//	return b
//}

func main() {
}
