package poker_test

import (
	"bytes"
	"fmt"
	poker "learning/17_HTTP"
	"strings"
	"testing"
)

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

			game := &poker.SpyGame{}
			cli := poker.NewCLI(game, in, stdout)

			err := cli.PlayPoker()

			poker.AssertNoError(t, err)

			assertMessagesSentToUser(t, stdout, poker.PlayerPrompt)
			poker.AssertStartGameNumberOfPlayers(t, game.StartCalledWith, test.StartPlayersInt)
			poker.AssertGameWinCalled(t, game, test.Name)
		})
	}

	t.Run("PlayPoker should return error when user inputs non number value as numberOfPlayers", func(t *testing.T) {
		in := strings.NewReader("u\n")
		stdout := &bytes.Buffer{}
		game := &poker.SpyGame{}
		cli := poker.NewCLI(game, in, stdout)

		err := cli.PlayPoker()

		if err != poker.InvalidInputError {
			t.Errorf("Expected error %v but got %v", poker.InvalidInputError, err)
		}

		assertMessagesSentToUser(t, stdout, poker.PlayerPrompt, poker.InvalidInput)
		assertGameNotStarted(t, game)
	})
}

func assertGameNotStarted(t *testing.T, game *poker.SpyGame) {
	t.Helper()

	if game.StartCalledWith != 0 {
		t.Fatalf("Game started when an error should have been thrown")
	}
}

func assertPlayerPrompt(t *testing.T, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func assertMessagesSentToUser(t *testing.T, stdout *bytes.Buffer, messages ...string) {
	t.Helper()
	want := strings.Join(messages, "")
	got := stdout.String()
	if got != want {
		t.Errorf("got %q sent to stdout but expected %+v", got, messages)
	}
}
