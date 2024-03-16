package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type User struct {
	Username   string
	LastActive string
}

type Database struct {
	Users map[string]User
}

func NewServer() *http.Server {
	return &http.Server{
		Addr: "localhost:9000",
	}
}

func initialDatabase() *Database {
	users := map[string]User{}

	for i := 1; i < 101; i++ {
		users[fmt.Sprintf("user%d", i)] = User{
			Username:   fmt.Sprintf("user%d", i),
			LastActive: time.Now().Format(time.RFC3339),
		}
	}

	return &Database{
		Users: users,
	}
}

func main() {
	server := NewServer()
	db := initialDatabase()

	fmt.Printf("Server listening on %s\n", server.Addr)

	http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		jsonString, err := json.Marshal(db.Users)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonString)
	})
	http.HandleFunc("/user", func(w http.ResponseWriter, r *http.Request) {
		username := r.URL.Query().Get("username")
		user, ok := db.Users[username]
		if !ok {
			http.Error(w, "user not found", http.StatusNotFound)
			return
		}

		jsonString, err := json.Marshal(user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonString)
	})

	server.ListenAndServe()
}
