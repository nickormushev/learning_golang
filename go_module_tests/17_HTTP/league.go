package poker

import (
	"encoding/json"
	"fmt"
	"io"
)

//League is an assortment of players
type League []Player

//NewLeague parses the read league into a []Player object
func NewLeague(read io.Reader) ([]Player, error) {
	var got []Player
	err := json.NewDecoder(read).Decode(&got)

	if err != nil {
		err = fmt.Errorf("Unable to parse response from server %q into slice of Player, '%v'", read, err)
	}

	return got, err
}

//Find is used to find the player entry throug a name
func (l League) Find(name string) *Player {
	for index, player := range l {
		if player.Name == name {
			return &l[index]
		}
	}

	return nil
}
