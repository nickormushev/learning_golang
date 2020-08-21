package poker

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

type StubPlayerStore struct {
	scores   map[string]int
	winCalls []string
	league   League
}

func (s StubPlayerStore) GetPlayerScore(playerName string) int {
	return s.scores[playerName]
}

func (s *StubPlayerStore) RecordWin(playerName string) {
	s.winCalls = append(s.winCalls, playerName)
}

func (s StubPlayerStore) GetLeague() League {
	return s.league
}

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
			request := newGetScoreRequest(test.username)
			response := httptest.NewRecorder()

			server.ServeHTTP(response, request)

			got := response.Body.String()

			assertResponseBody(t, got, test.expectedPoints)
			assertStatusCode(t, response.Code, http.StatusOK)
		})
	}

	t.Run("404 when player is missing", func(t *testing.T) {
		request := newGetScoreRequest("missing")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatusCode(t, response.Code, http.StatusNotFound)
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
		request := newPostWinRequest(player)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatusCode(t, response.Code, http.StatusAccepted)

		if len(store.winCalls) != 1 {
			t.Errorf("Win was not called")
		}

		if store.winCalls[0] != player {
			t.Errorf("did not store correct winner got %q want %q", store.winCalls[0], player)
		}
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

		request := newLeagueRequest()
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertJSONContentType(t, response)

		got := getLeagueFromResponse(t, response)

		assertLeague(t, got, wantedPlayers)
		assertStatusCode(t, response.Code, http.StatusOK)
	})
}

func assertSuccessfulParse(t *testing.T, err error, responseBody *bytes.Buffer) {

}

func assertStatusCode(t *testing.T, gotStatusCode, expectedStatusCode int) {
	t.Helper()

	if gotStatusCode != expectedStatusCode {
		t.Errorf("Status code mismatch: expected %d but got %d", expectedStatusCode, gotStatusCode)
	}
}

func assertLeague(t *testing.T, got, wantedLeague []Player) {
	t.Helper()
	if !reflect.DeepEqual(got, wantedLeague) {
		t.Fatalf("Wrong players were returned! Got: %v but wanted %v", got, wantedLeague)
	}
}

func newPostWinRequest(player string) *http.Request {
	request, _ := http.NewRequest(http.MethodPost, "/players/"+player, nil)
	return request
}

func newGetScoreRequest(player string) *http.Request {
	request, _ := http.NewRequest(http.MethodGet, "/players/"+player, nil)

	return request
}

func assertJSONContentType(t *testing.T, response *httptest.ResponseRecorder) {
	t.Helper()
	if response.Result().Header.Get("content-type") != jsonContentType {
		t.Errorf("response did not have content-type of application/json, got %v", response.Result().Header)
	}
}

func assertResponseBody(t *testing.T, responseBody, expectedResponseBody string) {
	if responseBody != expectedResponseBody {
		t.Errorf("Wanted %s but got %s", expectedResponseBody, responseBody)
	}
}

func newLeagueRequest() *http.Request {
	request, _ := http.NewRequest(http.MethodGet, "/league/", nil)
	return request
}

func getLeagueFromResponse(t *testing.T, resp *httptest.ResponseRecorder) []Player {
	got, err := NewLeague(resp.Body)

	if err != nil {
		t.Fatal(err)
	}

	return got
}
