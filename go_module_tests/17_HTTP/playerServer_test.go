package poker

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

var (
	dummyGame = &Game{}
	tenMS     = time.Millisecond * 10
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

	server := CreateNewPlayerServer(t, &store, dummyGame)

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

	server := CreateNewPlayerServer(t, &store, dummyGame)

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
		server := CreateNewPlayerServer(t, &store, dummyGame)

		request := NewLeagueRequest()
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		AssertJSONContentType(t, response)

		got := GetLeagueFromResponse(t, response)

		AssertLeague(t, got, wantedPlayers)
		AssertStatusCode(t, response.Code, http.StatusOK)
	})
}

func TestGame(t *testing.T) {
	t.Run("Should return 200 on get", func(t *testing.T) {
		store := &StubPlayerStore{}

		playerServer := CreateNewPlayerServer(t, store, dummyGame)

		request := newGameRequest()
		response := httptest.NewRecorder()

		playerServer.ServeHTTP(response, request)

		AssertStatusCode(t, response.Code, http.StatusOK)
	})

	t.Run("When we get a message over the websocket we declare it a winner of the game", func(t *testing.T) {
		wantedBlindAlert := "Blind is 100"
		game := &SpyGame{BlindAlert: []byte(wantedBlindAlert)}
		store := &StubPlayerStore{}
		playerServer := CreateNewPlayerServer(t, store, game)
		server := httptest.NewServer(playerServer)
		numberOfPlayers := "5"
		winner := "Jacob"

		ws := createWebSocket(t, "ws"+strings.TrimPrefix(server.URL, "http")+"/ws/")

		defer server.Close()
		defer ws.Close()

		sendWebSocketMessage(t, ws, numberOfPlayers)
		sendWebSocketMessage(t, ws, winner)

		AssertGameStartedWithXNumberOfPlayers(t, game, 5)
		AssertGameWinCalled(t, game, winner)
		within(t, tenMS, func() { assertWebsocketGotMsg(t, ws, wantedBlindAlert) })
	})
}

func newGameRequest() *http.Request {
	request, _ := http.NewRequest(http.MethodGet, "/game/", nil)
	return request
}

func assertWebsocketGotMsg(t *testing.T, ws *websocket.Conn, want string) {
	_, msg, _ := ws.ReadMessage()
	if string(msg) != want {
		t.Errorf(`got "%s", want "%s"`, string(msg), want)
	}
}

func createWebSocket(t *testing.T, wsURL string) *websocket.Conn {
	t.Helper()

	ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)

	if err != nil {
		t.Fatalf("could not open a ws connection on %s %v", wsURL, err)
	}

	return ws
}

func sendWebSocketMessage(t *testing.T, ws *websocket.Conn, msg string) {
	t.Helper()

	if err := ws.WriteMessage(websocket.TextMessage, []byte(msg)); err != nil {
		t.Fatalf("could not send message over ws connection %v", err)
	}
}

func within(t *testing.T, d time.Duration, assert func()) {
	t.Helper()

	done := make(chan struct{}, 1)

	go func() {
		assert()
		done <- struct{}{}
	}()

	select {
	case <-time.After(d):
		t.Error("timed out")
	case <-done:
	}
}
