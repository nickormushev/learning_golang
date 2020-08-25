package poker

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
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

func CreateTempFile(t *testing.T, initialData string) (*os.File, func()) {
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

func AssertNoError(t *testing.T, err error) {
	t.Helper()

	if err != nil {
		t.Fatalf("Unexpected error %v", err)
	}
}
