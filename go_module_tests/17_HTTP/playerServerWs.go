package poker

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var wsUpgrader websocket.Upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type playerServerWS struct {
	*websocket.Conn
}

func newPlayerServerWs(resp http.ResponseWriter, req *http.Request) *playerServerWS {
	conn, err := wsUpgrader.Upgrade(resp, req, nil)

	if err != nil {
		log.Printf("Problem upgrading http connection to web socket %v", err)
	}

	return &playerServerWS{conn}
}

func (p *playerServerWS) Write(msg []byte) (n int, err error) {
	err = p.WriteMessage(websocket.TextMessage, msg)

	if err != nil {
		return
	}

	return len(msg), nil
}

func (p *playerServerWS) WaitForMsg() string {
	_, msg, err := p.ReadMessage()

	if err != nil {
		log.Printf("Error reading msg %v", err)
	}

	return string(msg)
}
