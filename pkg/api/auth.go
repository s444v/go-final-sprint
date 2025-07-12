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
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, map[string]string{"error": "Invalid request"})
		return
	}
	correctPassword := os.Getenv("TODO_PASSWORD")
	if correctPassword == "" {
		w.WriteHeader(http.StatusInternalServerError)
		writeJSON(w, map[string]string{"error": "Server misconfiguration"})
		return
	}
	if pass.Password != correctPassword {
		w.WriteHeader(http.StatusUnauthorized)
		writeJSON(w, map[string]string{"error": "Unauthorized"})
		return
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": time.Now().Add(time.Hour * 2).Unix(),
		"sub": "todo-user",
	})

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		writeJSON(w, map[string]string{"error": err.Error()})
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
				w.WriteHeader(http.StatusUnauthorized)
				writeJSON(w, map[string]string{"error": "Authentication required"})
				return

			}
		}
		next(w, r)
	})
}
