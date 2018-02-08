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

	http.HandleFunc("/user-service", rootHandler)
	http.ListenAndServe(":3000", nil)
}
