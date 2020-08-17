package main

import (
	"context"
	"fmt"
	"net/http"
)

//Store is a object that is used to communicate with a store of data
type Store interface {
	Fetch(context.Context) (string, error)
	Cancel()
}

//Server takes in a Store(data) and returns it to the persone who requested it
func Server(store Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		data, err := store.Fetch(ctx)

		if err != nil {
			return
		}

		fmt.Fprintf(w, data)

		//Old implementation
		//data := make(chan string, 1)

		//go func() {
		//	data <- store.Fetch(ctx)
		//}()

		//select {
		//case res := <-data:
		//	fmt.Fprint(w, res)
		//case <-ctx.Done():
		//	close(data)
		//	store.Cancel()
		//}

	}
}

func main() {
	fmt.Println("vim-go")
}
