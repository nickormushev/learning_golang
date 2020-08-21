package poker

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
)

//FileSystemPlayerStore stores the player data in files
type FileSystemPlayerStore struct {
	database *json.Encoder
	league   League
}

//NewFileSystemPlayerStore is a constructor for FileSystemPlayer store that reads the database file
func NewFileSystemPlayerStore(file *os.File) (*FileSystemPlayerStore, error) {

	err := initialiseDbFile(file)

	if err != nil {
		return nil, fmt.Errorf("Could not initialise playerDb file %s %v", file.Name(), err)
	}

	league, err := NewLeague(file)

	if err != nil {
		return nil, fmt.Errorf("Failed to load player storef with this file %s %v", file.Name(), err)
	}

	return &FileSystemPlayerStore{
		database: json.NewEncoder(&tape{file}),
		league:   league,
	}, nil
}

func initialiseDbFile(file *os.File) error {
	file.Seek(0, 0)
	info, err := file.Stat()

	if err != nil {
		return fmt.Errorf("Failed to get file.Stat information for file %s %v", file.Name(), err)
	}

	if info.Size() == 0 {
		file.Write([]byte("[]"))
		file.Seek(0, 0)
	}

	return nil
}

//GetLeague reads the league from a file
func (f FileSystemPlayerStore) GetLeague() League {
	sort.Slice(f.league, func(fst, snd int) bool {
		return f.league[fst].Wins > f.league[snd].Wins
	})

	return f.league
}

//GetPlayerScore takes in a player name and returns their score
func (f FileSystemPlayerStore) GetPlayerScore(name string) int {
	player := f.league.Find(name)

	if player == nil {
		return 0
	}

	return player.Wins
}

//RecordWin updates a players win count
func (f *FileSystemPlayerStore) RecordWin(name string) {
	player := f.league.Find(name)

	if player != nil {
		player.Wins++
	} else {
		f.league = append(f.league, Player{name, 1})
	}

	f.database.Encode(f.league)
}
