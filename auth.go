package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/joho/godotenv"
)

type UserCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type TokenResponse struct {
	Token string `json:"token"`
}

type User struct {
	Username string
	Password string
}

// UserStore represents a simple in-memory user database.
var UserStore = map[string]User{
	"user1": {Username: "user1", Password: "pass1"},
	"user2": {Username: "user2", Password: "pass2"},
}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
}

// GenerateToken creates a JWT token for authenticated users.
func GenerateToken(user User) (string, error) {
	expirationTime := time.Now().Add(5 * time.Minute)
	claims := &jwt.StandardClaims{
		ExpiresAt: expirationTime.Unix(),
		Subject:   user.Username,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secretKey := os.Getenv("JWT_SECRET")
	if secretKey == "" {
		return "", fmt.Errorf("JWT_SECRET environment variable is not set")
	}

	return token.SignedString([]byte(secretKey))
}

// HandleUserAuthentication processes login requests and issues JWT tokens.
func HandleUserAuthentication(w http.ResponseWriter, r *http.Request) {
	var creds UserCredentials
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	user, validCredentials := UserStore[creds.Username]
	if !validCredentials || user.Password != creds.Password {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	token, err := GenerateToken(user)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(TokenResponse{Token: token})
}

// RequireTokenMiddleware ensures that the requester is authenticated.
func RequireTokenMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization token is required", http.StatusForbidden)
			return
		}

		tokenSplit := strings.Split(authHeader, " ")
		if len(tokenSplit) != 2 {
			http.Error(w, "Invalid or malformed authorization token", http.StatusForbidden)
			return
		}

		tokenPart := tokenSplit[1]

		claims := &jwt.StandardClaims{}
		token, err := jwt.ParseWithClaims(tokenPart, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Invalid authentication token", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// ProtectedEndpoint demonstrates a protected API endpoint.
func ProtectedEndpoint(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "You've accessed a protected endpoint!")
}

func main() {
	http.HandleFunc("/authenticate", HandleUserAuthentication)
	http.Handle("/protected", RequireTokenMiddleware(http.HandlerFunc(ProtectedEndpoint)))

	log.Println("Server is listening on port :8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Unable to start server: %v", err)
	}
}