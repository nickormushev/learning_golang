package main

import (
	"fmt"
	"io"
	"net/http"
)

//Greet takes in a name and says hello to the name
func Greet(writer io.Writer, name string) {
	fmt.Fprintf(writer, "Hello, %s!\n", name)
}

//MyGreeterHandler handles greet requests
func MyGreeterHandler(w http.ResponseWriter, r *http.Request) {
	Greet(w, "Niki")
}

func main() {
	http.ListenAndServe(":8000", http.HandlerFunc(MyGreeterHandler))
}
