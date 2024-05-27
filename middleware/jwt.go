package middleware

import (
    "net/http"
    "strings"

    "github.com/dgrijalva/jwt-go"
)

// JWTMiddleware is a middleware function to authenticate JWT tokens
func JWTMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Get the token from the Authorization header
        authHeader := r.Header.Get("Authorization")
        if authHeader == "" {
            http.Error(w, "Authorization header missing", http.StatusUnauthorized)
            return
        }

        // Check if the Authorization header has the format "Bearer <token>"
        parts := strings.Split(authHeader, " ")
        if len(parts) != 2 || parts[0] != "Bearer" {
            http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
            return
        }

        // Parse the JWT token
        tokenString := parts[1]
        token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
            // You should replace this with your own JWT secret key or key pair
            return []byte("your-secret-key"), nil
        })
        if err != nil {
            http.Error(w, "Failed to parse JWT token", http.StatusUnauthorized)
            return
        }

        // Check if the token is valid
        if !token.Valid {
            http.Error(w, "Invalid JWT token", http.StatusUnauthorized)
            return
        }

        // Call the next handler in the chain
        next.ServeHTTP(w, r)
    })
}
