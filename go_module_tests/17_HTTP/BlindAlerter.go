package poker

import (
	"fmt"
	"io"
	"time"
)

//BlindAlerter is an interface that represents alerters for blind values in a poker game
type BlindAlerter interface {
	ScheduledAlertAt(duration time.Duration, amount int, to io.Writer)
}

//BlindAlerterFunc is an implementation of BlindAlerter
type BlindAlerterFunc func(duration time.Duration, amount int, to io.Writer)

//ScheduledAlertAt takes in a duratiion and amount and calls a BlindAlerterFunc
func (a BlindAlerterFunc) ScheduledAlertAt(duration time.Duration, amount int, to io.Writer) {
	a(duration, amount, to)
}

//GenericAlerter writes a blind alert to the give io.Writer
func GenericAlerter(duration time.Duration, amount int, to io.Writer) {
	time.AfterFunc(duration, func() {
		fmt.Fprintf(to, "Blind is now %d\n", amount)
	})
}
