package poker_test

import (
	"bytes"
	"fmt"
	poker "learning/17_HTTP"
	"strings"
	"testing"
)

type SpyGame struct {
	startNumberOfPlayers int
	winCalls             []string
}

func (s *SpyGame) Start(numberOfPlayers int) {
	s.startNumberOfPlayers = numberOfPlayers
}

func (s *SpyGame) Win(winner string) {
	s.winCalls = append(s.winCalls, winner)
}

func TestCLI(t *testing.T) {
	cases := []struct {
		Name               string
		Input              string
		StartPlayersString string
		StartPlayersInt    int
	}{
		{Name: "Chris", Input: "Chris wins\n", StartPlayersString: "7\n", StartPlayersInt: 7},
		{Name: "Cleo", Input: "Cleo wins\n", StartPlayersString: "10\n", StartPlayersInt: 10},
	}

	for _, test := range cases {
		t.Run(fmt.Sprintf("%s wins", test.Name), func(t *testing.T) {
			in := strings.NewReader(test.StartPlayersString + test.Input)
			stdout := &bytes.Buffer{}

			game := &SpyGame{}
			cli := poker.NewCLI(game, in, stdout)

			err := cli.PlayPoker()

			poker.AssertNoError(t, err)

			assertMessagesSentToUser(t, stdout, poker.PlayerPrompt)
			assertStartGameNumberOfPlayers(t, game.startNumberOfPlayers, test.StartPlayersInt)
			assertGameWinCalled(t, game, test.Name)
		})
	}

	t.Run("PlayPoker should return error when user inputs non number value as numberOfPlayers", func(t *testing.T) {
		in := strings.NewReader("u\n")
		stdout := &bytes.Buffer{}
		game := &SpyGame{}
		cli := poker.NewCLI(game, in, stdout)

		err := cli.PlayPoker()

		if err != poker.InvalidInputError {
			t.Errorf("Expected error %v but got %v", poker.InvalidInputError, err)
		}

		assertMessagesSentToUser(t, stdout, poker.PlayerPrompt, poker.InvalidInput)
		assertGameNotStarted(t, game)
	})
}

func assertGameNotStarted(t *testing.T, game *SpyGame) {
	t.Helper()

	if game.startNumberOfPlayers != 0 {
		t.Fatalf("Game started when an error should have been thrown")
	}
}

func assertStartGameNumberOfPlayers(t *testing.T, got, want int) {
	t.Helper()

	if got != want {
		t.Errorf("got number of players %q, but wanted %q", got, want)
	}
}

func assertPlayerPrompt(t *testing.T, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func assertGameWinCalled(t *testing.T, game *SpyGame, player string) {
	t.Helper()
	poker.AssertElementInArray(t, game.winCalls, player)
}

func assertMessagesSentToUser(t *testing.T, stdout *bytes.Buffer, messages ...string) {
	t.Helper()
	want := strings.Join(messages, "")
	got := stdout.String()
	if got != want {
		t.Errorf("got %q sent to stdout but expected %+v", got, messages)
	}
}
