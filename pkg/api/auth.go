package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
)

func signinHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusBadRequest)
		err := writeJSON(w, map[string]string{"error": "wrong method"})
		if err != nil {
			http.Error(w, "cant parse to json", http.StatusInternalServerError)
		}
		return
	}
	var pass = struct {
		Password string `json:"password"`
	}{}
	err := json.NewDecoder(r.Body).Decode(&pass)
	if err != nil || pass.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		err = writeJSON(w, map[string]string{"error": "Invalid request"})
		if err != nil {
			http.Error(w, "cant parse to json", http.StatusInternalServerError)
		}
		return
	}
	correctPassword := TODO_PASSWORD
	if correctPassword == "" {
		w.WriteHeader(http.StatusInternalServerError)
		err = writeJSON(w, map[string]string{"error": "Server misconfiguration"})
		if err != nil {
			http.Error(w, "cant parse to json", http.StatusInternalServerError)
		}
		return
	}
	if pass.Password != correctPassword {
		w.WriteHeader(http.StatusUnauthorized)
		err = writeJSON(w, map[string]string{"error": "Unauthorized"})
		if err != nil {
			http.Error(w, "cant parse to json", http.StatusInternalServerError)
		}
		return
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": time.Now().Add(time.Hour * 2).Unix(),
		"sub": "todo-user",
	})

	tokenString, err := token.SignedString([]byte(JWT_SECRET))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		err = writeJSON(w, map[string]string{"error": err.Error()})
		if err != nil {
			http.Error(w, "cant parse to json", http.StatusInternalServerError)
		}
		return
	}
	w.Header().Set("Content-Type", "application/json")
	err = writeJSON(w, map[string]string{"token": tokenString})
	if err != nil {
		http.Error(w, "cant parse to json", http.StatusInternalServerError)
	}
}

func auth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// смотрим наличие пароля
		pass := TODO_PASSWORD
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
				return []byte(JWT_SECRET), nil
			})

			if err == nil && token.Valid {
				valid = true
			}

			if !valid {
				w.WriteHeader(http.StatusUnauthorized)
				err = writeJSON(w, map[string]string{"error": "Authentication required"})
				if err != nil {
					http.Error(w, "cant parse to json", http.StatusInternalServerError)
				}
				return
			}
		}
		next(w, r)
	})
}
