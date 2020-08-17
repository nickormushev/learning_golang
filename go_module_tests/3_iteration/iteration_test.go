package iteration

import (
	"fmt"
	"testing"
)

func TestIterate(t *testing.T) {
	got := Iterate("n", 5)
	want := "nnnnn"

	if want != got {
		t.Errorf("Expect %q and got %q", want, got)
	}
}

func BenchmarkIterate(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Iterate("n", 5)
	}
}

func ExampleIterate() {
	res := Iterate("n", 5)
	fmt.Println(res)
	// Output: nnnnn
}
