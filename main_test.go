package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"testing"
)

func TestGetUser(t *testing.T) {
	res, err := http.Get("http://localhost:9000/users")
	if err != nil {
		log.Fatal(err)
	}
	greeting, err := io.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s", greeting)
}
