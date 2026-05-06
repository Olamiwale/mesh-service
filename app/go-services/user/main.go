package main

import (
    "encoding/json"
    "log"
    "net/http"
    "os"
    "sync"
)

type User struct {
    ID    string `json:"id"`
    Email string `json:"email"`
    Name  string `json:"name"`
}

var users = map[string]User{
    "1": {ID: "1", Email: "raul@example.com", Name: "Raul Oguns"},
    "2": {ID: "2", Email: "paul@example.com", Name: "Paul Smith"},
    "3": {ID: "3", Email: "john@example.com", Name: "John Smith"},
    "4": {ID: "4", Email: "jane@example.com", Name: "Jane Smith"},
}

var mu sync.RWMutex

func main() {
    http.HandleFunc("/users/", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")

        id := r.URL.Path[len("/users/"):]

        mu.RLock()
        user, exists := users[id]
        mu.RUnlock()

        if !exists {
            http.Error(w, "User not found", 404)
            return
        }

        json.NewEncoder(w).Encode(user)
    })

    http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
        json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
    })

    // -------- FIXED PORT ISSUE --------
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    log.Println("User service running on port", port)
    http.ListenAndServe(":"+port, nil)
}
