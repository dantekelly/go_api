package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"testing"
)

func TestGetUser(t *testing.T) {
	t.Run("responds 100 users", func(t *testing.T) {
		res, err := http.Get("http://localhost:9000/users")
		if err != nil {
			log.Fatal(err)
		}

		if res.Header.Get("Content-Type") != "application/json" {
			t.Errorf("expected content type to be application/json, got %s", res.Header.Get("Content-Type"))
		}

		body, err := io.ReadAll(res.Body)
		res.Body.Close()

		if err != nil {
			log.Fatal(err)
		}

		var users map[string]User
		jsonErr := json.Unmarshal(body, &users)
		if jsonErr != nil {
			log.Fatal(jsonErr)
		}

		if len(users) != 100 {
			t.Errorf("expected 100 users, got %d", len(users))
		}

		fmt.Printf("%+v", users)
	})
}
