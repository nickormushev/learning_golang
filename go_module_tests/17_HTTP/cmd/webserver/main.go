package main

import (
	poker "learning/17_HTTP"
	"log"
	"net/http"
)

const dbFileName string = "game.db.json"

func main() {
	store, dbClose, err := poker.GenerateFileSystemPlayerStore(dbFileName)

	defer dbClose()

	if err != nil {
		log.Fatalf("Could not generate FileSystem player store from file, %v", err)
	}

	game := poker.NewGame(store, poker.BlindAlerterFunc(poker.GenericAlerter))
	server, err := poker.NewPlayerServer(store, game)

	if err != nil {
		log.Fatalf("Failed to create a playerServer %v", err)
	}

	//srv := &http.Server{Addr: ":5000", Handler: server}

	//go func() {
	//	if err = srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
	//		log.Fatalf("An erro occured %v", err)
	//	}
	//}()

	if err := http.ListenAndServe(":5000", server); err != nil {
		log.Fatalf("An erro occured %v", err)
	}
}
