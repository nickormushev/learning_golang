package poker

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"
	"time"
)

const (
	FullFileName   string = "./temp_config.yaml"
	FileName       string = "temp_config"
	TestDbFileName string = "test.db.json"
	TestServerPort string = "5000"
)

type SpyGame struct {
	StartCalled     bool
	StartCalledWith int
	BlindAlert      []byte

	WinCalled     bool
	WinCalledWith string
}

func (s *SpyGame) Start(numberOfPlayers int, to io.Writer) {
	s.StartCalled = true
	s.StartCalledWith = numberOfPlayers

	to.Write(s.BlindAlert)
}

func (s *SpyGame) Win(winner string) {
	s.WinCalled = true
	s.WinCalledWith = winner
}

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

func AssertTrueWithRetry(t *testing.T, got *bool) {
	t.Helper()

	passed := retryUntil(500*time.Millisecond, func() bool {
		return *got
	})

	if !passed {
		t.Errorf("expected true but got false")
	}
}

func AssertFalse(t *testing.T, got bool) {
	t.Helper()
	if got {
		t.Errorf("Expected false but got true")
	}
}

func AssertError(t *testing.T, err error) {
	t.Helper()

	if err == nil {
		t.Errorf("Expected error but got nil")
	}
}

func retryUntil(d time.Duration, f func() bool) bool {
	deadline := time.Now().Add(d)
	for time.Now().Before(deadline) {
		if f() {
			return true
		}
	}
	return false
}

func CreateNewPlayerServer(t *testing.T, store PlayerStore, game AbstractGame) *PlayerServer {
	t.Helper()
	server, err := NewPlayerServer(store, game)

	if err != nil {
		t.Fatalf("Failed to create NewPlayerServer with error %v", err)
	}

	return server
}

func AssertGameStartedWithXNumberOfPlayers(t *testing.T, game *SpyGame, numberOfPlayers int) {
	t.Helper()

	passed := retryUntil(500*time.Millisecond, func() bool {
		return game.StartCalled && game.StartCalledWith == numberOfPlayers
	})

	if !passed {
		t.Errorf("expected start called with %d but got %d", numberOfPlayers, game.StartCalledWith)
	}
}

func AssertStartGameNumberOfPlayers(t *testing.T, got, want int) {
	t.Helper()

	if got != want {
		t.Errorf("got number of players %q, but wanted %q", got, want)
	}
}

func AssertGameWinCalled(t *testing.T, game *SpyGame, player string) {
	t.Helper()

	passed := retryUntil(500*time.Millisecond, func() bool {
		return game.WinCalled && game.WinCalledWith == player
	})

	if !passed {
		t.Errorf("expected finish called with %q but got %q", player, game.WinCalledWith)
	}
}

func AssertElementInArray(t *testing.T, sl []string, want string) {
	t.Helper()

	if len(sl) != 1 {
		t.Errorf("Win was not called")
	}

	if sl[0] != want {
		t.Errorf("did not store correct winner got %q want %q", sl[0], want)
	}

}

func AssertUpdateWin(t *testing.T, store StubPlayerStore, player string) {
	t.Helper()
	AssertElementInArray(t, store.winCalls, player)
}

func AssertStatusCode(t *testing.T, gotStatusCode, expectedStatusCode int) {
	t.Helper()

	if gotStatusCode != expectedStatusCode {
		t.Errorf("Status code mismatch: expected %d but got %d", expectedStatusCode, gotStatusCode)
	}
}

func AssertLeague(t *testing.T, got, wantedLeague []Player) {
	t.Helper()
	if !reflect.DeepEqual(got, wantedLeague) {
		t.Fatalf("Wrong players were returned! Got: %v but wanted %v", got, wantedLeague)
	}
}

func NewPostWinRequest(player string) *http.Request {
	request, _ := http.NewRequest(http.MethodPost, "/players/"+player, nil)
	return request
}

func NewGetScoreRequest(player string) *http.Request {
	request, _ := http.NewRequest(http.MethodGet, "/players/"+player, nil)

	return request
}

func AssertJSONContentType(t *testing.T, response *httptest.ResponseRecorder) {
	t.Helper()
	if response.Result().Header.Get("content-type") != jsonContentType {
		t.Errorf("response did not have content-type of application/json, got %v", response.Result().Header)
	}
}

func AssertResponseBody(t *testing.T, responseBody, expectedResponseBody string) {
	if responseBody != expectedResponseBody {
		t.Errorf("Wanted %s but got %s", expectedResponseBody, responseBody)
	}
}

func NewLeagueRequest() *http.Request {
	request, _ := http.NewRequest(http.MethodGet, "/league/", nil)
	return request
}

func GetLeagueFromResponse(t *testing.T, resp *httptest.ResponseRecorder) []Player {
	got, err := NewLeague(resp.Body)

	if err != nil {
		t.Fatal(err)
	}

	return got
}

func AssertPlayerScore(t *testing.T, got, expectedPoints int) {
	if got != expectedPoints {
		t.Errorf("Should have gotten: %d points but he got %d", expectedPoints, got)
	}
}

//CreateTempFile creates a ioutil.temoFile  and returns it with a cancel() method to delete it
func CreateTempFile(t *testing.T, initialData, fileName string) (*os.File, func()) {
	t.Helper()
	tmpFile, err := ioutil.TempFile("", "db")

	return writeToNewlyCreatedFile(t, tmpFile, err, fileName, initialData)
}

//CreateTempFileOsOpenFile creates a normal file using os.OpenFile with a cancel() method to delete it
func CreateTempFileOsOpenFile(t *testing.T, initialData, filePath string) (*os.File, func()) {
	t.Helper()

	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0600)

	return writeToNewlyCreatedFile(t, file, err, filePath, initialData)
}

func writeToNewlyCreatedFile(t *testing.T, file *os.File, err error,
	fileName, initialData string) (*os.File, func()) {

	if err != nil {
		t.Fatalf("Could not create and open tempe file %v", err)
	}

	close := func() {
		file.Close()
		os.Remove(fileName)
	}

	_, err = file.WriteString(initialData)

	if err != nil {
		t.Fatalf("Failed to write to file %v", err)
	}

	file.Seek(0, 0)

	return file, close
}

//AssertNoError checks that the error given to it is nil and returns an error if it is not
func AssertNoError(t *testing.T, err error) {
	t.Helper()

	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}
}
