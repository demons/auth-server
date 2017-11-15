package main

import (
	"database/sql"
	"log"
	"net/http"

	"audiolang.com/auth-server/store"
	"audiolang.com/auth-server/tokgen"

	"audiolang.com/auth-server/models"
	"audiolang.com/auth-server/oauth"
	"audiolang.com/auth-server/utils"
)

// TODO: функции grant type должны прописать найденого пользователя в ctx и передать запрос дальше,

// GetGrantTypes возвратит реализацию аутентификации
// func GetGrantTypes(name string) func(ctx context.Context, w http.ResponseWriter, r *http.Request) (*models.User, bool) {
// 	switch name {
// 	case "password":
// 		{
// 			return grantTypePassword
// 		}
// 	case "code":
// 		{
// 			return grantTypeCode
// 		}
// 	case "refresh":
// 		{
// 			return grantTypeRefresh
// 		}
// 	}
// 	return nil
// }

// Реализации аутентификации

// GrantTypePassword аутентификация пользователя по паролю
func grantTypePassword(w http.ResponseWriter, r *http.Request) (*models.User, bool) {
	// Извлекаем из тела запроса email и password
	email := r.FormValue("email")
	password := r.FormValue("password")

	if email == "" {
		log.Println("Email is required")
		http.Error(w, "Email is required", http.StatusBadRequest)
		return nil, false
	}

	if password == "" {
		log.Printf("Password is required: %s\n", email)
		http.Error(w, "Password is required", http.StatusBadRequest)
		return nil, false
	}

	ctx := r.Context()

	// Вытаскиваем user store из ctx
	userStore, ok := store.FromContextWithUserStore(ctx)
	if ok == false {
		log.Printf("User store is not found")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return nil, false
	}

	// Получаем информацию о пользователе из базы данных
	user, err := userStore.FindByEmail(email)
	if err != nil {
		// Произошла какая-то ошибка при поиске
		log.Printf("Error finding a user: %v\n", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return nil, false
	}

	// Пользователь не найден
	if user == nil {
		log.Printf("User is not found: %v\n", err)
		http.Error(w, "Incorrect login or password", http.StatusBadRequest)
		return nil, false
	}

	// Проверяем пароль пользователя
	err = utils.VerifyPassword(user.Hash.String, user.Salt.String, password)
	if err != nil {
		log.Printf("Incorrect password: %v\n", err)
		http.Error(w, "Incorrect login or password", http.StatusBadRequest)
		return nil, false
	}

	return user, true
}

// GrantTypeCode аутентификация пользователя по коду, который вернул сервер соц. сети
func grantTypeCode(w http.ResponseWriter, r *http.Request) (*models.User, bool) {
	// Найти указанный провайдер
	// Обменять code на accessToken
	// Получить профиль пользователя
	// Найти этого пользователя в базе данных, если не найден, то создать
	providerName := r.FormValue("provider")
	code := r.FormValue("code")

	if providerName == "" {
		log.Println("Provider name is required")
		http.Error(w, "Provider name is required", http.StatusBadRequest)
		return nil, false
	}

	if code == "" {
		log.Println("Code is required")
		http.Error(w, "Code is required", http.StatusBadRequest)
		return nil, false
	}

	ctx := r.Context()

	// Вытаскиваем user store из ctx
	userStore, ok := store.FromContextWithUserStore(ctx)
	if ok == false {
		log.Printf("User store is not found")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return nil, false
	}

	// Ищем провайдера по названию
	provider, err := oauth.GetProviderByName(providerName)
	if err != nil {
		log.Printf("Error getting provider by provider name: %v", err)
		http.Error(w, "This provider is not supported", http.StatusBadRequest)
		return nil, false
	}

	// Обмениваем code на access token
	accessToken, err := provider.ExchangeCode(code)
	if err != nil {
		log.Printf("Error exchanging code: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return nil, false
	}

	// Получаем профиль пользователя
	userProfile, err := provider.GetUserProfile(accessToken)
	if err != nil {
		log.Printf("Error getting a user profile: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return nil, false
	}

	sid := userProfile.ProviderName + ":" + userProfile.ID

	// Ищем пользователя в базе данных по SID (social id = providerName:ID)
	user, err := userStore.FindBySID(sid)
	if err != nil {
		// Пользователь не найден, нужно создать нового (не факт конечно, что ошибка вызвана из-за того, что пользователь не найден, а не из-за недоступности базы данных) ПОДУМАЙ
		newUser := models.User{
			SID:      sql.NullString{String: sid, Valid: true},
			Name:     sql.NullString{String: userProfile.Name, Valid: true},
			IsSocial: true,
			IsActive: true,
		}
		userID, err := userStore.Insert(&newUser)
		if err != nil {
			log.Printf("Error inserting a new user: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return nil, false
		}
		newUser.ID = userID

		user = &newUser
	}

	return user, true
}

// GrantTypeRefresh аутентификация по refresh токену
func grantTypeRefresh(w http.ResponseWriter, r *http.Request) (*models.User, *models.Token, bool) {
	refresh := r.FormValue("refresh")

	if refresh == "" {
		log.Println("Refresh is required")
		http.Error(w, "Refresh is required", http.StatusBadRequest)
		return nil, nil, false
	}

	ctx := r.Context()

	// Вытаскиваем user store из ctx
	userStore, ok := store.FromContextWithUserStore(ctx)
	if ok == false {
		log.Printf("User store is not found")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return nil, nil, false
	}

	// Вытаскиваем tokenGenerator из ctx
	tokenGenerator, ok := tokgen.FromContextWithTokenGenerator(ctx)
	if ok == false {
		log.Printf("TokenGenerator is not found")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return nil, nil, false
	}

	// Ищмем указанный токен
	token, err := tokenGenerator.FindToken(ctx, refresh)
	if err != nil {
		log.Printf("Error finding token: %v, token: %s", err, refresh)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return nil, nil, false
	}

	// Если токен не найден
	if token == nil {
		log.Printf("Token is not found: %v, token: %s", err, refresh)
		http.Error(w, "Access denied", http.StatusForbidden)
		return nil, nil, false
	}

	// Проверяем валидность токена
	if token.Valid() == false {
		log.Printf("Token is not valid: token: %s", refresh)
		http.Error(w, "Access denied", http.StatusForbidden)
		return nil, nil, false
	}

	// Получаем информацию о пользователе из базы данных
	user, err := userStore.FindByUserID(token.UserID)
	if err != nil {
		log.Printf("Error finding user: %v\n", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return nil, nil, false
	}

	if user == nil {
		log.Printf("User is not found: %v\n", err)
		http.Error(w, "Access denied", http.StatusForbidden)
		return nil, nil, false
	}

	// Обновляем токен
	newToken, err := tokenGenerator.UpdateToken(ctx, token)
	if err != nil {
		log.Printf("Error updating refresh token: %v token: %s", err, refresh)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return nil, nil, false
	}

	return user, newToken, true
}
