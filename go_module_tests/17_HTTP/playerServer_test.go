package poker

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGETPlayerScore(t *testing.T) {

	store := StubPlayerStore{
		map[string]int{
			"gosho": 20,
			"pesho": 10,
		},
		nil,
		nil,
	}

	server := NewPlayerServer(&store)

	cases := []struct {
		username       string
		expectedPoints string
	}{
		{"gosho", "20"},
		{"pesho", "10"},
	}

	for _, test := range cases {
		t.Run("Valid get"+test.username, func(t *testing.T) {
			request := NewGetScoreRequest(test.username)
			response := httptest.NewRecorder()

			server.ServeHTTP(response, request)

			got := response.Body.String()

			AssertResponseBody(t, got, test.expectedPoints)
			AssertStatusCode(t, response.Code, http.StatusOK)
		})
	}

	t.Run("404 when player is missing", func(t *testing.T) {
		request := NewGetScoreRequest("missing")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		AssertStatusCode(t, response.Code, http.StatusNotFound)
	})
}

func TestStoreWins(t *testing.T) {
	store := StubPlayerStore{
		map[string]int{},
		nil,
		nil,
	}
	server := NewPlayerServer(&store)

	t.Run("Records when I send a POST request", func(t *testing.T) {
		player := "gosho"
		request := NewPostWinRequest(player)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		AssertStatusCode(t, response.Code, http.StatusAccepted)

		AssertUpdateWin(t, store, player)
	})
}

func TestLeague(t *testing.T) {
	t.Run("Return 200 on /league", func(t *testing.T) {
		wantedPlayers := League{
			{"Chris", 21},
			{"George", 53},
			{"Doki", 20},
		}

		store := StubPlayerStore{league: wantedPlayers}
		server := NewPlayerServer(&store)

		request := NewLeagueRequest()
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		AssertJSONContentType(t, response)

		got := GetLeagueFromResponse(t, response)

		AssertLeague(t, got, wantedPlayers)
		AssertStatusCode(t, response.Code, http.StatusOK)
	})
}
