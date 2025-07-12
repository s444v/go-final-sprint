package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

var jwtSecret = []byte("12345")

func signinHandler(w http.ResponseWriter, r *http.Request) {
	var pass = struct {
		Password string `json:"password"`
	}{}
	err := json.NewDecoder(r.Body).Decode(&pass)
	if err != nil || pass.Password == "" {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	correctPassword := os.Getenv("TODO_PASSWORD")
	if correctPassword == "" {
		http.Error(w, "Server misconfiguration", http.StatusInternalServerError)
		return
	}
	if pass.Password != correctPassword {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": time.Now().Add(time.Hour * 2).Unix(),
		"sub": "todo-user",
	})

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"token": tokenString,
	})

}

func auth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// смотрим наличие пароля
		pass := os.Getenv("TODO_PASSWORD")
		if len(pass) > 0 {
			var jwtToken string // JWT-токен из куки
			// получаем куку
			cookie, err := r.Cookie("token")
			if err == nil {
				jwtToken = cookie.Value
			}

			var valid bool
			// Валидация JWT
			token, err := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
				// Проверка типа алгоритма
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				return jwtSecret, nil
			})

			if err == nil && token.Valid {
				valid = true
			}

			if !valid {
				http.Error(w, "Authentication required", http.StatusUnauthorized)
				return
			}
		}
		next(w, r)
	})
}
