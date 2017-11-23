package main

import (
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

// HandlePasswordReset восстановление пароля
func HandlePasswordReset(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	email := r.PostFormValue("email")

	if email == "" {
		log.Println("Email is empty")
		http.Error(w, "Email is required", http.StatusBadRequest)
		return
	}

	// Проверяем формат указанного email
	err := utils.VerifyEmailFormat(email)
	if err != nil {
		log.Printf("Invalid email format: %v", err)
		http.Error(w, "Invalid email fromat", http.StatusBadRequest)
		return
	}

	// TODO: Нужно ограничить количество отправленных подрад писем

	// Выполняем поиск пользователя по указанному email
	findedUser, err := userDb.FindByEmail(email)
	if err != nil {
		log.Printf("Error finding user: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if findedUser == nil {
		// Такого пользователя не существует
		log.Printf("User is not found")
		http.Error(w, "User is not found", http.StatusBadRequest)
		return
	}

	// Проверим тип аккаунта
	if findedUser.IsSocial == true {
		// Нельзя восстановить пароль к аккаунту, созданному с помощью соц. сетей
		log.Println("Password reset failed. Account is not native")
		http.Error(w, "Server internal error", http.StatusInternalServerError)
		return
	}

	ctx := r.Context()

	// Устанавливаем в context пользователея
	ctx = models.NewContextWithUser(ctx, findedUser)

	// Устанавливаем в context хранилище временных токенов
	ctx = store.NewContextWithTokenStore(ctx, tempTokenDb)

	// Генерируем временный токен доступа
	token, err := tempTokenGenerator.CreateToken(ctx, []string{"password_reset"})
	if err != nil {
		log.Printf("Error creating a temp token: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Установим токену права на сброс пароля
	token.SetScope("reset_password")

	// Отправляем по почте код активации для подтверждения аккаунта
	template := messageTemplates["resetPassword"]
	go emailNotificator.SendResetPasswordMessage(template, findedUser.Email.String, token.Token)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("A letter has been sent to the e-mail to reset the password"))
}

// HandlePasswordChange смена пароля
func HandlePasswordChange(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Извлеч пользователя из контекста
	ctx := r.Context()

	user, ok := models.FromContextWithUser(ctx)
	if ok != true {
		log.Printf("Error reading user from context")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Считать старый пароль и новый
	oldPassword := r.FormValue("oldPassword")
	if oldPassword == "" {
		log.Println("OldPassword is empty")
		http.Error(w, "OldPassword is required", http.StatusBadRequest)
		return
	}

	newPassword := r.FormValue("newPassword")
	if newPassword == "" {
		log.Println("NewPassword is empty")
		http.Error(w, "NewPassword is required", http.StatusBadRequest)
		return
	}

	// Проверяем верный ли старый пароль
	err := utils.VerifyPassword(user.Hash.String, user.Salt.String, oldPassword)
	if err != nil {
		log.Printf("Incorrect password: %v\n", err)
		http.Error(w, "Incorrect login or password", http.StatusBadRequest)
		return
	}

	// Изменить пароль

	// Создаем хэш пароля
	hash, salt, err := utils.HashPasswordWithSalt(newPassword)
	if err != nil {
		log.Printf("Error creating hash password: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	user.Hash = sql.NullString{String: hash, Valid: true}
	user.Salt = sql.NullString{String: salt, Valid: true}

	// Создаем нового пользователя
	err = userDb.UpdatePassword(user)
	if err != nil {
		log.Printf("Error creating a user: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Password changed successfully"))
}

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

	ctx := r.Context()

	// Устанавливаем в context пользователея
	ctx = models.NewContextWithUser(ctx, &newUser)

	// Устанавливаем в context хранилище временных токенов
	ctx = store.NewContextWithTokenStore(ctx, tempTokenDb)

	r = r.WithContext(ctx)

	// Установим токену права на активацию аккаунта
	token, err := tempTokenGenerator.CreateToken(ctx, []string{"activate_account"})
	if err != nil {
		log.Printf("Error creating a temp token: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Отправляем по почте код активации для подтверждения аккаунта
	template := messageTemplates["activateAccount"]
	go emailNotificator.SendActivationCode(template, newUser.Email.String, token.Token)

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("User registered successfully"))
}

// HandleToken returns access token
func HandleToken(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

	// start := time.Now()

	grantType := r.FormValue("grant_type")
	log.Println("get token with grant_type:", grantType)

	ctx := r.Context()

	// Устанавливаем в context хранилище пользователей
	ctx = store.NewContextWithUserStore(ctx, userDb)

	// Устанавливаем в context хранилище refresh токенов
	ctx = store.NewContextWithTokenStore(ctx, refreshTokenDb)

	// Устанавливаем в context генератор токенов
	ctx = tokgen.NewContextWithTokenGenerator(ctx, tokenGenerator)

	r = r.WithContext(ctx)

	// Инициализируем нового пользователя
	var user *models.User
	var refreshToken *models.Token
	var ok bool

	switch grantType {
	case "password":
		user, ok = grantTypePassword(w, r)
	case "code":
		user, ok = grantTypeCode(w, r)
	case "refresh":
		user, refreshToken, ok = grantTypeRefresh(w, r)

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

	r = r.WithContext(ctx)

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
		newRefreshToken, err := tokenGenerator.CreateToken(ctx, []string{"refresh_token"})
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
