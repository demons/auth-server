package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"audiolang.com/auth-server/tokgen"

	"audiolang.com/auth-server/utils"

	"audiolang.com/auth-server/models"
	"audiolang.com/auth-server/store"
	"github.com/julienschmidt/httprouter"
)

// HandleSignUp регистрация нового пользователя
func HandleSignUp(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	email := r.FormValue("email")
	password := r.FormValue("password")

	if email == "" {
		log.Println("Email is empty")
		http.Error(w, "Email is required", http.StatusBadRequest)
		return
	}

	if password == "" {
		log.Println("Password is empty")
		http.Error(w, "Password is required", http.StatusBadRequest)
		return
	}

	// Проверяем email на корректность
	err := utils.VerifyEmailFormat(email)
	if err != nil {
		log.Printf("Invalid email format: %v", err)
		http.Error(w, "Invalid email format", http.StatusBadRequest)
		return
	}

	// Проверяем существует ли пользователь с таким username
	searchUser, err := userDb.FindByEmail(email)
	if err != nil {
		log.Printf("Error finding user: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if searchUser != nil {
		// Такой пользователь уже существует
		log.Println("User already exists")
		http.Error(w, "User already exists", http.StatusBadRequest)
		return
	}

	// Создаем хэш пароля
	hash, salt, err := utils.HashPasswordWithSalt(password)
	if err != nil {
		log.Printf("Error creating hash password: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	newUser := models.User{
		Email:    sql.NullString{String: email, Valid: true},
		Hash:     sql.NullString{String: hash, Valid: true},
		Salt:     sql.NullString{String: salt, Valid: true},
		IsActive: true,
		IsSocial: false,
	}

	// Создаем нового пользователя
	userID, err := userDb.Insert(&newUser)
	if err != nil {
		log.Printf("Error creating a user: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	newUser.ID = userID

	ctx := context.Background()

	// Устанавливаем в context пользователея
	ctx = models.NewContextWithUser(ctx, &newUser)

	// Устанавливаем в context хранилище временных токенов
	ctx = store.NewContextWithTokenStore(ctx, tempTokenDb)

	token, err := tempTokenGenerator.CreateToken(ctx)
	if err != nil {
		log.Printf("Error creating a temp token: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// TODO: Нужно отправить письмо с подтверждением email, пользователю на почту

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("User registered successfully. " + token.Token))
}

// HandleToken returns access token
func HandleToken(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

	// start := time.Now()

	grantType := r.FormValue("grant_type")
	log.Println("get token with grant_type:", grantType)

	ctx := context.Background()

	// Устанавливаем в context хранилище пользователей
	ctx = store.NewContextWithUserStore(ctx, userDb)

	// Устанавливаем в context хранилище refresh токенов
	ctx = store.NewContextWithTokenStore(ctx, refreshTokenDb)

	// Устанавливаем в context генератор токенов
	ctx = tokgen.NewContextWithTokenGenerator(ctx, tokenGenerator)

	// Инициализируем нового пользователя
	var user *models.User
	var refreshToken *models.Token
	var ok bool

	switch grantType {
	case "password":
		user, ok = grantTypePassword(ctx, w, r)
	case "code":
		user, ok = grantTypeCode(ctx, w, r)
	case "refresh":
		user, refreshToken, ok = grantTypeRefresh(ctx, w, r)
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
	if refreshToken == nil {
		// Генерируем новый refresh token
		newRefreshToken, err := tokenGenerator.CreateToken(ctx)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		data.RefreshToken = newRefreshToken.Token
		data.RefreshExpires = newRefreshToken.Expires
	} else {
		data.RefreshToken = refreshToken.Token
		data.RefreshExpires = refreshToken.Expires
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
