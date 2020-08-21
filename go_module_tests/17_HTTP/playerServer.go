package poker

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

const jsonContentType string = "application/json"

//PlayerStore contains the information of the players
type PlayerStore interface {
	GetPlayerScore(string) int
	RecordWin(string)
	GetLeague() League
}

//PlayerServer is the httpHandler for request to /players/
type PlayerServer struct {
	store PlayerStore
	http.Handler
}

//Player represents a person with a name and a number of wins
type Player struct {
	Name string
	Wins int
}

//NewPlayerServer is a constructor for PlayerServer that creates a router for it.
func NewPlayerServer(store PlayerStore) *PlayerServer {
	p := new(PlayerServer)

	p.store = store

	router := http.NewServeMux()
	router.Handle("/players/", http.HandlerFunc(p.playersHandler))
	router.Handle("/league/", http.HandlerFunc(p.leagueHandler))

	p.Handler = router

	return p
}

func (p *PlayerServer) leagueHandler(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("content-type", jsonContentType)
	json.NewEncoder(resp).Encode(p.store.GetLeague())
}

func (p *PlayerServer) playersHandler(resp http.ResponseWriter, req *http.Request) {
	player := strings.TrimPrefix(req.URL.Path, "/players/")

	switch req.Method {
	case http.MethodPost:
		p.processWin(resp, player)
	case http.MethodGet:
		p.displayScore(resp, player)
	}
}

func (p *PlayerServer) processWin(resp http.ResponseWriter, player string) {
	p.store.RecordWin(player)
	resp.WriteHeader(http.StatusAccepted)
}

func (p *PlayerServer) displayScore(resp http.ResponseWriter, player string) {

	score := p.store.GetPlayerScore(player)

	if score == 0 {
		resp.WriteHeader(http.StatusNotFound)
	}

	fmt.Fprint(resp, score)
}
