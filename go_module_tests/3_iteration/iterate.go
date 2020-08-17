package iteration

//Iterate takes in a character and a count and repeats characters count times
func Iterate(c string, count int) string {
	var res string

	for i := 0; i < count; i++ {
		res += c
	}

	return res
}
