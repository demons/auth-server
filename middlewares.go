package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"audiolang.com/auth-server/models"
	"audiolang.com/auth-server/tokgen"
	"github.com/julienschmidt/httprouter"
)

// Metric измерение продолжительности выполнения запроса
func Metric(f func(w http.ResponseWriter, r *http.Request, params httprouter.Params)) func(http.ResponseWriter, *http.Request, httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		timeStart := int(time.Now().UnixNano())
		f(w, r, params)
		// Переведем наносекунды в милисекунды
		ms := (int(time.Now().UnixNano()) - timeStart) / 1000000
		log.Printf("duration: %d ms, url: %s", ms, r.URL.Path)
	}
}

// Logger middleware
func Logger(f func(w http.ResponseWriter, r *http.Request, params httprouter.Params)) func(http.ResponseWriter, *http.Request, httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		log.Println("Hello, I'm logger")
		f(w, r, params)
	}
}

// Auth with temp token or authorization token
func Auth(f func(w http.ResponseWriter, r *http.Request, params httprouter.Params)) func(http.ResponseWriter, *http.Request, httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		tempToken := r.FormValue("code")
		authorization := r.Header.Get("Authorization")

		var user *models.User
		var err error

		if tempToken != "" {
			user, err = authWithTempToken(tempToken)
		} else if authorization != "" {
			authLines := strings.SplitN(authorization, " ", 2)
			if len(authLines) != 2 || strings.ToLower(authLines[0]) != "bearer" {
				log.Println("Authorization header format must be Bearer {token}")
				http.Error(w, "Authorization header format must be Bearer {token}", http.StatusUnauthorized)
				return
			}
			tokenString := authLines[1]
			user, err = authWithToken(tokenString)
		}

		if err != nil {
			log.Println(err)
			http.Error(w, "Authentication error", http.StatusUnauthorized)
			return
		}

		if user == nil {
			http.Error(w, "Authentication error", http.StatusUnauthorized)
			return
		}

		ctx := r.Context()

		// Set user to context
		ctx = models.NewContextWithUser(ctx, user)

		r = r.WithContext(ctx)

		f(w, r, params)

	}
}

// AuthWithTempToken auth with temp token or authorization token
func authWithTempToken(tempToken string) (*models.User, error) {

	// Найти временный токен в хранилище
	token, err := tempTokenDb.FindByToken(tempToken)
	if err != nil {
		return nil, fmt.Errorf("Error finding token: %v", err)
	}

	if token == nil {
		return nil, fmt.Errorf("Token is not found")
	}

	if token.Valid() == false {
		return nil, fmt.Errorf("Token is not valid")
	}

	// Find user from store by temp token
	user, err := userDb.FindByUserID(token.UserID)
	if err != nil {
		return nil, fmt.Errorf("Error finding user: %v", err)
	}

	if user == nil {
		return nil, fmt.Errorf("User is not found")
	}

	return user, nil
}

// AuthWithToken is auth with authorization
func authWithToken(token string) (*models.User, error) {
	// Check token
	claims, err := tokenChecker.Check(token)
	if err != nil {
		if err == tokgen.ErrTokenIsInvalid {
			return nil, fmt.Errorf("Token is not valid")
		}
		return nil, fmt.Errorf("Error checking token: %v", err)
	}

	userID := int64(claims["userID"].(float64))

	// Find user from store by userID
	user, err := userDb.FindByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("Error finding user: %v", err)
	}

	if user == nil {
		return nil, fmt.Errorf("User is not found")
	}

	return user, nil

}
