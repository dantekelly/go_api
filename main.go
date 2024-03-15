package main

import (
	"fmt"
	"net/http"
)

func NewServer() *http.Server {
	return &http.Server{
		Addr: "localhost:9000",
	}
}

func main() {
	server := NewServer()

	fmt.Printf("Server listening on %s\n", server.Addr)

	http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	})

	server.ListenAndServe()
}
