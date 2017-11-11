package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"audiolang.com/auth-server/tokgen"

	"audiolang.com/auth-server/db"
	"audiolang.com/auth-server/models"
	"github.com/julienschmidt/httprouter"
)

// A very simple health check.
// w.WriteHeader(http.StatusOK)
// w.Header().Set("Content-Type", "application/json")

// In the future we could report back on the status of our DB, or our cache
// (e.g. Redis) by performing a simple PING, and include them in the response.
// io.WriteString(w, `{"alive": true}`)

// Коментарии. Можно написать функцию loginWithPasswordHandler(w http.WriterResponse, login, password string)

// HandleToken returns access token
func HandleToken(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

	// start := time.Now()

	grantType := r.FormValue("grant_type")
	log.Println("get token with grant_type:", grantType)

	ctx := context.Background()

	// Устанавливаем в context хранилище пользователей
	ctx = db.NewContextWithUserStore(ctx, userStore)

	// Устанавливаем в context хранилище refresh токенов
	ctx = db.NewContextWithRefreshTokenStore(ctx, refreshTokenStore)

	// Инициализируем нового пользователя
	var user *models.User
	var updatedRefresh *models.RefreshToken
	var ok bool

	switch grantType {
	case "password":
		user, ok = grantTypePassword(ctx, w, r)
	case "code":
		user, ok = grantTypeCode(ctx, w, r)
	case "refresh":
		user, updatedRefresh, ok = grantTypeRefresh(ctx, w, r)
	default:
		{
			log.Println("This grant type is not supported")
			http.Error(w, "This grant type is not supported", http.StatusBadRequest)
			return
		}
	}

	if ok == false {
		return
	}

	// Устанавливаем в context пользователя
	ctx = models.NewContextWithUser(ctx, user)

	// Генерируем токен доступа для пользователя
	accessToken, err := jwtGen.Token(ctx)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	var data struct {
		AccessToken    string `json:"access_token"`
		RefreshToken   string `json:"refresh_token"`
		RefreshExpires int64  `json:"refresh_exp"`
	}
	data.AccessToken = accessToken

	// Если refresh token не существует (grant_type != refresh), то генерируем новый
	if updatedRefresh == nil {
		// Генерируем новый refresh token
		refresh, err := tokgen.NewRefreshToken(ctx)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		data.RefreshToken = refresh.Token
		data.RefreshExpires = refresh.Expires
	} else {
		data.RefreshToken = updatedRefresh.Token
		data.RefreshExpires = updatedRefresh.Expires
	}

	// Формируем ответ
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(&data)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// elapsed := time.Since(start)
	// log.Println(elapsed)
}
