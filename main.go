package main

import (
	"fmt"
	"net/http"
)

func rootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello world")
}

func main() {
	fmt.Printf("hello world")

	http.HandleFunc("/", rootHandler)
	http.ListenAndServe(":8080", nil)
}
