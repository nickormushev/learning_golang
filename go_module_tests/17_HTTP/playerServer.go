package poker

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"
)

const (
	jsonContentType string = "application/json"
	gamePath        string = "./html/game.html"
)

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
	tmpl *template.Template
	game AbstractGame
}

//Player represents a person with a name and a number of wins
type Player struct {
	Name string
	Wins int
}

//NewPlayerServer is a constructor for PlayerServer that creates a router for it.
func NewPlayerServer(store PlayerStore, game AbstractGame) (*PlayerServer, error) {
	p := new(PlayerServer)

	tmpl, err := template.ParseFiles(gamePath)

	if err != nil {
		return nil, fmt.Errorf("Error loading template %v", err)
	}

	p.tmpl = tmpl
	p.store = store
	p.game = game

	router := http.NewServeMux()
	router.Handle("/players/", http.HandlerFunc(p.playersHandler))
	router.Handle("/league/", http.HandlerFunc(p.leagueHandler))
	router.Handle("/game/", http.HandlerFunc(p.gameHandler))
	router.Handle("/ws/", http.HandlerFunc(p.webSocketHandler))

	p.Handler = router

	return p, nil
}

func (p *PlayerServer) webSocketHandler(resp http.ResponseWriter, req *http.Request) {
	conn := newPlayerServerWs(resp, req)
	numberOfPlayers, _ := strconv.Atoi(conn.WaitForMsg())

	p.game.Start(numberOfPlayers, conn)

	winnerMsg := conn.WaitForMsg()
	p.game.Win(string(winnerMsg))
}

func (p *PlayerServer) gameHandler(resp http.ResponseWriter, req *http.Request) {
	p.tmpl.Execute(resp, nil)
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
