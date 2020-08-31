package poker

import (
	"fmt"
	"testing"
)

const fileName string = "db"

func TestGETLeague(t *testing.T) {
	t.Run("/league from a reader", func(t *testing.T) {
		database, cleanDb := CreateTempFile(t, `[
            {"Name": "Cleo", "Wins": 10},
            {"Name": "Chris", "Wins": 33}]`, fileName)

		defer cleanDb()

		store, err := NewFileSystemPlayerStore(database)
		AssertNoError(t, err)

		got := store.GetLeague()

		want := League{
			{Name: "Chris", Wins: 33},
			{Name: "Cleo", Wins: 10},
		}

		got = store.GetLeague()
		AssertLeague(t, got, want)
	})

	t.Run("/league returns sorted slice", func(t *testing.T) {
		database, cleanDb := CreateTempFile(t, `[
            {"Name": "Cleo", "Wins": 10},
            {"Name": "Chris", "Wins": 33},
				{"Name": "Joro", "Wins": 12},
				{"Name": "Kiro", "Wins": 22}]`, fileName)

		defer cleanDb()

		store, err := NewFileSystemPlayerStore(database)
		AssertNoError(t, err)

		league := store.GetLeague()

		wanted := League{
			{Name: "Chris", Wins: 33},
			{Name: "Kiro", Wins: 22},
			{Name: "Joro", Wins: 12},
			{Name: "Cleo", Wins: 10},
		}

		AssertLeague(t, league, wanted)
	})
}

func TestPlayerScore(t *testing.T) {
	cases := []Player{
		{Name: "Cleo", Wins: 10},
		{Name: "Chris", Wins: 33},
	}

	for _, test := range cases {
		t.Run(fmt.Sprintf("/players/%s score from a reader", test.Name), func(t *testing.T) {
			database, cleanDb := CreateTempFile(t, `[
            {"Name": "Cleo", "Wins": 10},
            {"Name": "Chris", "Wins": 33}]`, fileName)

			defer cleanDb()

			store, err := NewFileSystemPlayerStore(database)
			AssertNoError(t, err)

			got := store.GetPlayerScore(test.Name)

			AssertPlayerScore(t, got, test.Wins)
		})
	}

	t.Run("Test update player score", func(t *testing.T) {
		database, cleanDb := CreateTempFile(t, `[
            {"Name": "Cleo", "Wins": 10},
            {"Name": "Chris", "Wins": 33}]`, fileName)

		defer cleanDb()

		store, err := NewFileSystemPlayerStore(database)
		AssertNoError(t, err)

		player := "Chris"
		store.RecordWin(player)
		expectedPoints := 34

		AssertPlayerScore(t, store.GetPlayerScore(player), expectedPoints)
	})

	t.Run("Update should create new user if non exists", func(t *testing.T) {
		database, cleanDb := CreateTempFile(t, `[
    	      {"Name": "Cleo", "Wins": 10},
    	      {"Name": "Chris", "Wins": 33}]`, fileName)

		defer cleanDb()

		store, err := NewFileSystemPlayerStore(database)
		AssertNoError(t, err)

		player := "Missing"
		store.RecordWin(player)
		expectedPoints := 1

		AssertPlayerScore(t, store.GetPlayerScore(player), expectedPoints)

	})
}

func TestWorksWithEmptyFiles(t *testing.T) {
	t.Run("works with an empty file", func(t *testing.T) {
		database, cleanDatabase := CreateTempFile(t, "", fileName)
		defer cleanDatabase()

		_, err := NewFileSystemPlayerStore(database)

		AssertNoError(t, err)
	})
}
