package poker

import (
	"fmt"
	"os"
)

//GenerateFileSystemPlayerStore generates reads a FileSystemPlayerStore from given file and returns it
func GenerateFileSystemPlayerStore(dbFileName string) (*FileSystemPlayerStore, func(), error) {
	file, err := os.OpenFile(dbFileName, os.O_RDWR|os.O_CREATE, 0666)

	if err != nil {
		return nil, nil, fmt.Errorf("Could not open file %s %v", dbFileName, err)
	}

	store, err := NewFileSystemPlayerStore(file)

	if err != nil {
		return nil, nil, fmt.Errorf("Could not create File System player store %v", err)
	}

	closeFunc := func() {
		file.Close()
	}

	return store, closeFunc, nil
}
