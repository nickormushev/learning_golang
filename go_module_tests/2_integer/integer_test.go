package integer

import (
	"fmt"
	"testing"
)

func TestAdder(t *testing.T) {
	t.Run("Add two numbers", func(t *testing.T) {
		got := Add(3, 4)
		expect := 7

		if got != expect {
			t.Errorf("The addition should give %d but is %d", expect, got)
		}
	})
}

func ExampleAdd() {
	sum := Add(1, 5)
	fmt.Println(sum)
	// Output: 6
}
