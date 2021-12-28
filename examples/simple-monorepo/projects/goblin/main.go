package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", func (w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Welcome to goblin! demo service👾")
	})
	println("Goblin demo service started on port 8000")
	http.ListenAndServe(":8000", nil)
}