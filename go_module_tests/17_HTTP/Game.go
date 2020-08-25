package poker

import (
	"io"
	"time"
)

type AbstractGame interface {
	Win(winner string)
	Start(numberOfPlayers int, to io.Writer)
}

type Game struct {
	store   PlayerStore
	alerter BlindAlerter
}

//NewGame is a constructor for Game
func NewGame(store PlayerStore, alerter BlindAlerter) AbstractGame {
	return &Game{store, alerter}
}

//Win takes in user input and records a winner
func (g *Game) Win(winner string) {
	g.store.RecordWin(winner)
}

//Start is the beggining of the game and it alerts the increase of the blind value after a set amount of time. Tha follows the formula 5 + numberOfPlayers = time.Minutes until increment
func (g *Game) Start(numberOfPlayers int, to io.Writer) {
	blindIncrement := time.Duration(5+numberOfPlayers) * time.Minute

	blinds := []int{100, 200, 300, 400, 500, 600, 800, 1000, 2000, 4000, 8000}
	blindTime := 0 * time.Minute
	for _, blind := range blinds {
		g.alerter.ScheduledAlertAt(blindTime, blind, to)
		blindTime = blindTime + blindIncrement
	}
}
