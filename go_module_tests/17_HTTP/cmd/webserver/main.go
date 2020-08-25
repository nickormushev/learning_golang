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

	server := poker.NewPlayerServer(store)
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
