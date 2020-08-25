package main

import (
	"fmt"
	poker "learning/17_HTTP"
	"log"
	"os"
)

var dbFileName string = "cli.db.json"

func main() {

	store, dbClose, err := poker.GenerateFileSystemPlayerStore(dbFileName)

	defer dbClose()

	if err != nil {
		log.Fatalf("Could not generate FileSystem player store from file, %v", err)
	}

	game := poker.NewGame(store, poker.BlindAlerterFunc(poker.StdOutAlerter))
	gameCLI := poker.NewCLI(game, os.Stdin, os.Stdout)

	fmt.Print("It's poker time\n")
	fmt.Println("Type {Name} wins to record a win")
	gameCLI.PlayPoker()
}
