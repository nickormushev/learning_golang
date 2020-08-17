package integer

import "testing"

func TestAdder(t *testing.T) {
	t.Run("Add two numbers", func(t *testing.T) {
		got := add(3, 4)
		expect := 7

		if got != expect {
			t.Errorf("The addition should give %q but is %q", expect, got)
		}
	})
}
