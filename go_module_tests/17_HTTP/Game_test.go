package poker_test

import (
	"fmt"
	poker "learning/17_HTTP"
	"testing"
	"time"
)

type SpyBlindAlerter struct {
	alerts []scheduledAlert
}

type scheduledAlert struct {
	at     time.Duration
	amount int
}

func (s scheduledAlert) String() string {
	return fmt.Sprintf("%d chips at %v", s.amount, s.at)
}

func (s *SpyBlindAlerter) ScheduledAlertAt(duration time.Duration, amount int) {
	s.alerts = append(s.alerts, scheduledAlert{duration, amount})
}

func TestRecordWin(t *testing.T) {
	cases := []struct {
		Name string
	}{
		{Name: "Chris"},
		{Name: "Cleo"},
	}

	var dummySpyAlerter = &SpyBlindAlerter{}

	for _, test := range cases {
		t.Run(fmt.Sprintf("%s wins", test.Name), func(t *testing.T) {
			playerStore := &poker.StubPlayerStore{}
			game := poker.NewGame(playerStore, dummySpyAlerter)

			game.Win(test.Name)

			poker.AssertUpdateWin(t, *playerStore, test.Name)
		})
	}
}

func TestStart(t *testing.T) {
	cases := []struct {
		numberOfPlayers int
		alerts          []scheduledAlert
	}{
		{
			numberOfPlayers: 7,
			alerts: []scheduledAlert{
				{0 * time.Minute, 100},
				{12 * time.Minute, 200},
				{24 * time.Minute, 300},
				{36 * time.Minute, 400},
				{48 * time.Minute, 500},
			},
		},
		{
			numberOfPlayers: 5,
			alerts: []scheduledAlert{
				{0 * time.Minute, 100},
				{10 * time.Minute, 200},
				{20 * time.Minute, 300},
				{30 * time.Minute, 400},
				{40 * time.Minute, 500},
				{50 * time.Minute, 600},
				{60 * time.Minute, 800},
				{70 * time.Minute, 1000},
				{80 * time.Minute, 2000},
				{90 * time.Minute, 4000},
				{100 * time.Minute, 8000},
			},
		},
	}

	for _, test := range cases {
		t.Run("Blind alerts are triggered after a certain amount of time", func(t *testing.T) {
			playerStore := &poker.StubPlayerStore{}
			blindAlerter := &SpyBlindAlerter{}

			game := poker.NewGame(playerStore, blindAlerter)
			game.Start(test.numberOfPlayers)

			for i, test := range test.alerts {
				t.Run(fmt.Sprintf("Amount %d at time %v", test.amount, test.at), func(t *testing.T) {
					if len(blindAlerter.alerts) <= i {
						t.Fatalf("Expected %d alerts but got %d", i, len(blindAlerter.alerts))
					}

					alert := blindAlerter.alerts[i]

					assertAlert(t, alert, test)
				})
			}
		})

	}
}

func assertAlert(t *testing.T, gotAlert, wantedAlert scheduledAlert) {
	t.Helper()

	if gotAlert.amount != wantedAlert.amount {
		t.Fatalf("Expected amount %d alerts but got %d", wantedAlert.at, gotAlert.amount)
	}

	if gotAlert.at != wantedAlert.at {
		t.Errorf("Expected time %d alerts but got %d", wantedAlert.at, gotAlert.at)
	}
}
