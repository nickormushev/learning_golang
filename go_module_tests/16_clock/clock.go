package clockface

import (
	"fmt"
	"time"
)

//Point represents a point in euclidian space
type Point struct {
	X int
	Y int
}

type Clock struct {
}

func SecondHand(t time.Time) Point {
	return Point{150, 60}
}

func main() {
	fmt.Println("vim-go")
}
