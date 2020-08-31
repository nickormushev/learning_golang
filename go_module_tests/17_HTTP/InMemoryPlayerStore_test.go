package poker

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRecordingWinsAndRetreivingThem(t *testing.T) {
	database, cleanDb := CreateTempFile(t, "[]", "db")
	defer cleanDb()

	store, err := NewFileSystemPlayerStore(database)
	AssertNoError(t, err)

	game := &Game{}
	server := CreateNewPlayerServer(t, store, game)

	player := "Pepper"

	server.ServeHTTP(httptest.NewRecorder(), NewPostWinRequest(player))
	server.ServeHTTP(httptest.NewRecorder(), NewPostWinRequest(player))
	server.ServeHTTP(httptest.NewRecorder(), NewPostWinRequest(player))

	t.Run("Get player scores", func(t *testing.T) {
		response := httptest.NewRecorder()
		server.ServeHTTP(response, NewGetScoreRequest(player))
		AssertStatusCode(t, response.Code, http.StatusOK)

		AssertResponseBody(t, response.Body.String(), "3")
	})

	t.Run("Get league", func(t *testing.T) {
		response := httptest.NewRecorder()
		server.ServeHTTP(response, NewLeagueRequest())

		AssertStatusCode(t, response.Code, http.StatusOK)

		wantedPlayers := []Player{
			{player, 3},
		}

		got := GetLeagueFromResponse(t, response)

		AssertLeague(t, got, wantedPlayers)
	})
}
