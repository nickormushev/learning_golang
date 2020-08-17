package arrays

import (
	"reflect"
	"testing"
)

func TestSum(t *testing.T) {
	t.Run("Run test with slice", func(t *testing.T) {
		arr := []int{1, 2, 3, 4, 5}
		got := Sum(arr)
		want := 15

		if got != want {
			t.Errorf("The sum of %v we got was %d but we expected %d", arr, got, want)
		}
	})
}

func TestSumAll(t *testing.T) {
	t.Run("With two arguments", func(t *testing.T) {
		want := []int{3, 9}
		got := SumAll(nil, []int{1, 2}, []int{0, 9})

		if !reflect.DeepEqual(got, want) {
			t.Errorf("We expected %v and got %v", want, got)
		}
	})

}

func TestSumAllTails(t *testing.T) {
	want := []int{6, 0, 29, 10}
	got := SumAllTails([]int{1, 2, 4}, []int{}, []int{0, 9, 20}, []int{5, 6, 4})

	if !reflect.DeepEqual(want, got) {
		t.Errorf("We expected %v but got %v", want, got)
	}
}
