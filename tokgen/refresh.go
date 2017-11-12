package tokgen

import (
	"context"
	"errors"
	"log"
	"time"

	"audiolang.com/auth-server/models"
	"audiolang.com/auth-server/store"
)

func createRefreshToken(userID int64) *models.RefreshToken {
	// Создаем новый refresh token
	return &models.RefreshToken{
		UserID:  userID,
		Token:   generateToken(),
		Expires: time.Now().Add(time.Hour * 24 * 10).Unix(), // Токен истекает через 10 дней
	}
}

// NewRefreshToken генерирует новый refresh token
func NewRefreshToken(ctx context.Context) (*models.RefreshToken, error) {
	// Извлекаем из контекста пользователя
	user, ok := models.FromContextWithUser(ctx)
	if ok == false {
		log.Println("User is not found in context")
		return nil, errors.New("User is not found in context")
	}

	// Извлекаем из контекста refreshTokenStore
	refreshTokenStore, ok := store.FromContextWithRefreshTokenStore(ctx)
	if ok == false {
		log.Println("Refresh store is not found in context")
		return nil, errors.New("Refresh store is not found in context")
	}

	// Создаем новый токен
	refreshToken := createRefreshToken(user.ID)

	// Добавляем токен в базу данных
	err := refreshTokenStore.Insert(refreshToken)
	if err != nil {
		log.Printf("Error inserting new refresh token: %v\n", err)
		return nil, errors.New("Error inserting new refresh token")
	}

	return refreshToken, nil
}

// ChangeRefreshToken обновляем имеющийся токен, возвращаем новый токен
func ChangeRefreshToken(ctx context.Context, token string) (*models.RefreshToken, error) {

	// Извлекаем из контекста refreshTokenStore
	refreshTokenStore, ok := store.FromContextWithRefreshTokenStore(ctx)
	if ok == false {
		log.Println("Refresh store is not found in context")
		return nil, errors.New("Refresh store is not found in context")
	}

	// Проверяем не истек ли токен
	oldRefreshToken, err := refreshTokenStore.FindByRefreshToken(token)
	if err != nil {
		// Произошла какая-то ошибка при поиске токена
		log.Printf("Error finding refresh token: %v\n", err)
		return nil, errors.New("Error finding refresh token")
	}

	if oldRefreshToken == nil {
		// Токен не найден в хранилище
		log.Println("Refresh token is not found")
		return nil, errors.New("Refresh token is not found")
	}

	now := time.Now().Unix()
	if now >= oldRefreshToken.Expires {
		// Токен истек
		log.Printf("Refresh token has expired: now: %d, token_exp: %d\n", now, oldRefreshToken.Expires)
		return nil, errors.New("Refresh token has expired")
	}

	// Создаем новый токен
	refreshToken := createRefreshToken(oldRefreshToken.UserID)

	// Обновляем старый, активный токен
	updatedRefresh, err := refreshTokenStore.Update(token, refreshToken)
	if err != nil {
		log.Printf("Error updating refresh token: %v\n", err)
		return nil, errors.New("Error updating refresh token")
	}

	return updatedRefresh, nil
}
