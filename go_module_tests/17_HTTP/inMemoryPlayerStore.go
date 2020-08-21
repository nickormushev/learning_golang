package poker

//This fail is redundant now and replaced by the FileSystemPlayerStore
import "sync"

//InMemoryPlayerStore is the in memory store for players
type InMemoryPlayerStore struct {
	scores map[string]int
	mx     sync.Mutex
}

//NewInMemoryPlayerStore is a constructor for the InMemoryPlayerStore
func NewInMemoryPlayerStore() *InMemoryPlayerStore {
	return &InMemoryPlayerStore{scores: map[string]int{}}
}

//GetPlayerScore takes in a player name and returns the score of that player
func (i *InMemoryPlayerStore) GetPlayerScore(name string) int {
	return i.scores[name]
}

//RecordWin takes in a name and increments the number of his wins
func (i *InMemoryPlayerStore) RecordWin(name string) {
	i.mx.Lock()
	i.scores[name]++

	i.mx.Unlock()
}

//GetLeague returns the all the players in the league
func (i *InMemoryPlayerStore) GetLeague() []Player {
	var players []Player

	for name, score := range i.scores {
		players = append(players, Player{name, score})
	}

	return players
}
