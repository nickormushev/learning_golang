package main

import (
	poker "learning/17_HTTP"
	server "learning/17_HTTP/server"
	"log"
)

const (
	dbFileName string = "game.db.json"
	serverPort string = ":5000"
)

func main() {
	store, dbClose, err := poker.GenerateFileSystemPlayerStore(dbFileName)

	defer dbClose()

	if err != nil {
		log.Fatalf("Could not generate FileSystem player store from file, %v", err)
	}

	game := poker.NewGame(store, poker.BlindAlerterFunc(poker.GenericAlerter))
	playerServer, err := poker.NewPlayerServer(store, game)

	ctx := server.GenerateContextWithSigint()
	srv := server.CreateServerAndServe(ctx, playerServer, serverPort)

	<-ctx.Done()
	err = server.GracefullShutdown(ctx, srv)

	if err != nil {
		log.Fatalf("Failed to shutdown gracefully %v!", err)
	} else {
		log.Printf("Successful graceful shutdown of server")
	}
}
