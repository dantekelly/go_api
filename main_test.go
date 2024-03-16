package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"testing"
)

func TestGetUsers(t *testing.T) {
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

func TestGetUser(t *testing.T) {
	t.Run("responds 1 user for good user", func(t *testing.T) {
		res, err := http.Get("http://localhost:9000/user?username=user1")
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

		var user User
		jsonErr := json.Unmarshal(body, &user)
		if jsonErr != nil {
			log.Fatal(jsonErr)
		}

		if user.Username != "user1" {
			t.Errorf("expected username to be user1, got %s", user.Username)
		}
	})

	t.Run("responds 404 for bad query", func(t *testing.T) {
		res, err := http.Get("http://localhost:9000/user?username=baduser")
		if err != nil {
			log.Fatal(err)
		}

		if res.StatusCode != http.StatusNotFound {
			t.Errorf("expected status code to be 404, got %d", res.StatusCode)
		}
	})
}
