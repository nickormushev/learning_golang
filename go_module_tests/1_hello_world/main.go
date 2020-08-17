package main

import "fmt"

//Hello says hello
func Hello(name string) string {
	if name == "" {
		return "Hello, world!\n"
	}
	return "Hello, " + name + "!\n"
}

func main() {
	fmt.Printf(Hello(""))
}
