package routes

import (
    "database/sql"
    "encoding/json"
    "net/http"
    "time"

    "golang.org/x/crypto/bcrypt"
)

// User represents a user in the database
type User struct {
    ID       int       `json:"id"`
    Username string    `json:"username"`
    Email    string    `json:"email"`
    Password string    `json:"password"`
    CreatedAt time.Time `json:"created_at"`
}

// SignupHandler handles user signup requests
func SignupHandler(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var user User
        if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
            http.Error(w, "Invalid input", http.StatusBadRequest)
            return
        }

        // Hash the password
        hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
        if err != nil {
            http.Error(w, "Error creating user", http.StatusInternalServerError)
            return
        }

        // Insert the new user into the database
        query := "INSERT INTO users (username, email, password) VALUES ($1, $2, $3) RETURNING id, created_at"
        err = db.QueryRow(query, user.Username, user.Email, string(hashedPassword)).Scan(&user.ID, &user.CreatedAt)
        if err != nil {
            http.Error(w, "Error creating user", http.StatusInternalServerError)
            return
        }

        w.WriteHeader(http.StatusCreated)
        json.NewEncoder(w).Encode(user)
    }
}

// LoginHandler handles user login requests
func LoginHandler(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var credentials struct {
            Username string `json:"username"`
            Password string `json:"password"`
        }
        if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
            http.Error(w, "Invalid input", http.StatusBadRequest)
            return
        }

        var storedUser User
        query := "SELECT id, username, email, password, created_at FROM users WHERE username=$1"
        err := db.QueryRow(query, credentials.Username).Scan(&storedUser.ID, &storedUser.Username, &storedUser.Email, &storedUser.Password, &storedUser.CreatedAt)
        if err != nil {
            if err == sql.ErrNoRows {
                http.Error(w, "Invalid username or password", http.StatusUnauthorized)
            } else {
                http.Error(w, "Error logging in", http.StatusInternalServerError)
            }
            return
        }

        // Compare the stored hashed password with the provided password
        if err := bcrypt.CompareHashAndPassword([]byte(storedUser.Password), []byte(credentials.Password)); err != nil {
            http.Error(w, "Invalid username or password", http.StatusUnauthorized)
            return
        }

        w.WriteHeader(http.StatusOK)
        json.NewEncoder(w).Encode(storedUser)
    }
}
