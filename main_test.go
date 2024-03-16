package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"testing"
)

var baseUrl = "http://localhost:9000"

func TestGetUsers(t *testing.T) {
	t.Run("responds 100 users", func(t *testing.T) {
		res, err := http.Get(fmt.Sprintf("%s/users", baseUrl))
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

func BenchmarkGetUser(b *testing.B) {
	for i := 0; i < b.N; i++ {
		randInt := rand.Intn(100-1+1) + 1
		res, err := http.Get(fmt.Sprintf("%s/user?username=user%d", baseUrl, randInt))
		if err != nil {
			log.Fatal(err)
		}

		if res.Header.Get("Content-Type") != "application/json" {
			b.Errorf("expected content type to be application/json, got %s", res.Header.Get("Content-Type"))
		}

		body, err := io.ReadAll(res.Body)
		res.Body.Close()

		if err != nil {
			log.Fatal(err)
		}

		var user User
		jsonErr := json.Unmarshal(body, &user)
		if jsonErr != nil {
			fmt.Println("error")
			log.Fatal(jsonErr)
		}

		if user.Username != fmt.Sprintf("user%d", randInt) {
			b.Errorf("expected username to be user1, got %s", user.Username)
		}
	}
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
