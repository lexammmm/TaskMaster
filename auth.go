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

	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func AuthenticateUser(w http.ResponseWriter, r *http.Request) {
	var userCredentials UserCredentials
	err := json.NewDecoder(r.Body).Decode(&userCredentials)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user, ok := users[userCredentials.Username]
	if !ok || user.Password != userCredentials.Password {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	token, err := GenerateJWT(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(JwtToken{Token: token})
}

func JwtAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenHeader := r.Header.Get("Authorization")
		if tokenHeader == "" {
			w.WriteHeader(http.StatusForbidden)
			fmt.Fprintln(w, "Missing auth token")
			return
		}

		splitted := strings.Split(tokenHeader, " ")
		if len(splitted) != 2 {
			w.WriteHeader(http.StatusForbidden)
			fmt.Fprintln(w, "Invalid/Malformed auth token")
			return
		}

		tokenPart := splitted[1] 

		claims := &jwt.StandardClaims{}
		token, err := jwt.ParseWithClaims(tokenPart, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil {
			w.WriteHeader(http.StatusForbidden)
			fmt.Fprintf(w, "Malformed authentication token: %v", err)
			return
		}

		if !token.Valid {
			w.WriteHeader(http.StatusForbidden)
			fmt.Fprintln(w, "Token is not valid.")
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
	http.ListenAndServe(":8080", nil)
}