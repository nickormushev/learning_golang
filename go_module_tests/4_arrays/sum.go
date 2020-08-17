package arrays

//Sum takes an array and returns the sum of it's elements
func Sum(arr []int) int {
	sum := 0

	for _, v := range arr {
		sum += v
	}

	return sum
}

//SumAll takes in a variable amount of arrays and sums each one separately
func SumAll(f func([]int) int, numbersToSum ...[]int) (sums []int) {
	if f == nil {
		f = Sum
	}

	for k := range numbersToSum {
		if len(numbersToSum[k]) == 0 {
			sums = append(sums, 0)
		} else {
			sums = append(sums, f(numbersToSum[k]))
		}
	}

	return
}

//SumAllTails takes in a variable amount of arrays and sums the tails of each array
func SumAllTails(tailsToSum ...[]int) []int {
	return SumAll(func(arr []int) int { return Sum(arr[1:]) }, tailsToSum...)
}
