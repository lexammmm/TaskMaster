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

type JwtToken struct {
	Token string `json:"token"`
}

type User struct {
	Username string
	Password string
}

var users = map[string]User{
	"user1": {Username: "user1", Password: "pass1"},
	"user2": {Username: "user2", Password: "pass2"},
}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
}

func GenerateJWT(user User) (string, error) {
	expirationTime := time.Now().Add(5 * time.Minute)
	claims := &jwt.StandardClaims{
		ExpiresAt: expirationTime.Unix(),
		Subject:   user.Username,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secretKey := os.Getenv("JWT_SECRET")
	if secretKey == "" {
		return "", fmt.Errorf("JWT_SECRET is not set")
	}

	return token.SignedString([]byte(secretKey))
}

func AuthenticateUser(w http.ResponseWriter, r *http.Request) {
	var userCredentials UserCredentials
	if err := json.NewDecoder(r.Body).Decode(&userCredentials); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	user, ok := users[userCredentials.Username]
	if !ok || user.Password != userCredentials.Password {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	token, err := GenerateJWT(user)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(JwtToken{Token: token})
}

func JwtAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenHeader := r.Header.Get("Authorization")
		if tokenHeader == "" {
			http.Error(w, "Missing auth token", http.StatusForbidden)
			return
		}

		splitToken := strings.Split(tokenHeader, " ")
		if len(splitToken) != 2 {
			http.Error(w, "Invalid/Malformed auth token", http.StatusForbidden)
			return
		}

		tokenPart := splitToken[1]

		claims := &jwt.StandardClaims{}
		token, err := jwt.ParseWithClaims(tokenPart, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Malformed authentication token", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func TestProtectedEndpoint(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Congratulations! This is a protected endpoint.")
}

func main() {
	http.HandleFunc("/authenticate", AuthenticateUser)
	http.Handle("/protected", JwtAuthentication(http.HandlerFunc(TestProtectedEndpoint)))

	log.Println("Listening on port :8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}