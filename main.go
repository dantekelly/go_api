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

type Cache struct {
	Users    map[string]*User
	Attempts int
}

type Database struct {
	Users    map[string]*User
	Attempts int
}

type Server struct {
	db    *Database
	cache *Cache
	serve *http.Server
}

func NewServer() *Server {
	users := make(map[string]*User)

	for i := 1; i < 101; i++ {
		users[fmt.Sprintf("user%d", i)] = &User{
			Username:   fmt.Sprintf("user%d", i),
			LastActive: time.Now().Format(time.RFC3339),
		}
	}

	return &Server{
		db:    &Database{Users: users, Attempts: 0},
		cache: &Cache{Users: make(map[string]*User), Attempts: 0},
		serve: &http.Server{
			Addr: "localhost:9000",
		},
	}
}

func (s *Server) tryCache(username string) (*User, error) {
	user, ok := s.cache.Users[username]
	s.cache.Attempts++
	if !ok {
		dbUser, dbOk := s.db.Users[username]
		s.db.Attempts++

		if !dbOk {
			return nil, fmt.Errorf("user not found")
		}

		s.cache.Users[username] = dbUser
		return dbUser, nil
	}

	return user, nil
}

func (s *Server) handleGetUsers(w http.ResponseWriter, r *http.Request) {
	jsonString, err := json.Marshal(s.db.Users)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonString)
}
func (s *Server) handleGetUser(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	user, err := s.tryCache(username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	jsonString, jsonErr := json.Marshal(user)
	if jsonErr != nil {
		http.Error(w, jsonErr.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonString)
}

func main() {
	s := NewServer()

	fmt.Printf("Server listening on %s\n", s.serve.Addr)

	http.HandleFunc("/users", s.handleGetUsers)
	http.HandleFunc("/user", s.handleGetUser)

	s.serve.ListenAndServe()
}
