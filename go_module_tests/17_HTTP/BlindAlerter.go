package poker

import (
	"fmt"
	"os"
	"time"
)

//BlindAlerter is an interface that represents alerters for blind values in a poker game
type BlindAlerter interface {
	ScheduledAlertAt(duration time.Duration, amount int)
}

//BlindAlerterFunc is an implementation of BlindAlerter
type BlindAlerterFunc func(duration time.Duration, amount int)

//ScheduledAlertAt takes in a duratiion and amount and calls a BlindAlerterFunc
func (a BlindAlerterFunc) ScheduledAlertAt(duration time.Duration, amount int) {
	a(duration, amount)
}

//StdOutAlerter writes a blind alert to stdout
func StdOutAlerter(duration time.Duration, amount int) {
	time.AfterFunc(duration, func() {
		fmt.Fprintf(os.Stdout, "Blind is now %d\n", amount)
	})
}
