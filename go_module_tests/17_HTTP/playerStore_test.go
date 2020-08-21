package poker

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

func TestGETLeague(t *testing.T) {
	t.Run("/league from a reader", func(t *testing.T) {
		database, cleanDb := createTempFile(t, `[
            {"Name": "Cleo", "Wins": 10},
            {"Name": "Chris", "Wins": 33}]`)

		defer cleanDb()

		store, err := NewFileSystemPlayerStore(database)
		assertNoError(t, err)

		got := store.GetLeague()

		want := League{
			{Name: "Chris", Wins: 33},
			{Name: "Cleo", Wins: 10},
		}

		got = store.GetLeague()
		assertLeague(t, got, want)
	})

	t.Run("/league returns sorted slice", func(t *testing.T) {
		database, cleanDb := createTempFile(t, `[
            {"Name": "Cleo", "Wins": 10},
            {"Name": "Chris", "Wins": 33},
				{"Name": "Joro", "Wins": 12},
				{"Name": "Kiro", "Wins": 22}]`)

		defer cleanDb()

		store, err := NewFileSystemPlayerStore(database)
		assertNoError(t, err)

		league := store.GetLeague()

		wanted := League{
			{Name: "Chris", Wins: 33},
			{Name: "Kiro", Wins: 22},
			{Name: "Joro", Wins: 12},
			{Name: "Cleo", Wins: 10},
		}

		assertLeague(t, league, wanted)
	})
}

func TestPlayerScore(t *testing.T) {
	cases := []Player{
		{Name: "Cleo", Wins: 10},
		{Name: "Chris", Wins: 33},
	}

	for _, test := range cases {
		t.Run(fmt.Sprintf("/players/%s score from a reader", test.Name), func(t *testing.T) {
			database, cleanDb := createTempFile(t, `[
            {"Name": "Cleo", "Wins": 10},
            {"Name": "Chris", "Wins": 33}]`)

			defer cleanDb()

			store, err := NewFileSystemPlayerStore(database)
			assertNoError(t, err)

			got := store.GetPlayerScore(test.Name)

			assertPlayerScore(t, got, test.Wins)
		})
	}

	t.Run("Test update player score", func(t *testing.T) {
		database, cleanDb := createTempFile(t, `[
            {"Name": "Cleo", "Wins": 10},
            {"Name": "Chris", "Wins": 33}]`)

		defer cleanDb()

		store, err := NewFileSystemPlayerStore(database)
		assertNoError(t, err)

		player := "Chris"
		store.RecordWin(player)
		expectedPoints := 34

		assertPlayerScore(t, store.GetPlayerScore(player), expectedPoints)
	})

	t.Run("Update should create new user if non exists", func(t *testing.T) {
		database, cleanDb := createTempFile(t, `[
    	      {"Name": "Cleo", "Wins": 10},
    	      {"Name": "Chris", "Wins": 33}]`)

		defer cleanDb()

		store, err := NewFileSystemPlayerStore(database)
		assertNoError(t, err)

		player := "Missing"
		store.RecordWin(player)
		expectedPoints := 1

		assertPlayerScore(t, store.GetPlayerScore(player), expectedPoints)

	})
}

func TestWorksWithEmptyFiles(t *testing.T) {
	t.Run("works with an empty file", func(t *testing.T) {
		database, cleanDatabase := createTempFile(t, "")
		defer cleanDatabase()

		_, err := NewFileSystemPlayerStore(database)

		assertNoError(t, err)
	})
}

func assertPlayerScore(t *testing.T, got, expectedPoints int) {
	if got != expectedPoints {
		t.Errorf("Should have gotten: %d points but he got %d", expectedPoints, got)
	}
}

func createTempFile(t *testing.T, initialData string) (*os.File, func()) {
	t.Helper()
	tmpFile, err := ioutil.TempFile("", "db")

	if err != nil {
		t.Fatalf("Could not open file %v", err)
	}

	tmpFile.WriteString(initialData)

	removeFile := func() {
		tmpFile.Close()
		os.Remove(tmpFile.Name())
	}

	return tmpFile, removeFile
}

func assertNoError(t *testing.T, err error) {
	t.Helper()

	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}
}
