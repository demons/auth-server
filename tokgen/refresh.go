package tokgen

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"audiolang.com/auth-server/models"
	"audiolang.com/auth-server/store"
)

func generateRefreshToken() *models.RefreshToken {
	// Последовательность случайных байт
	buf := make([]byte, 128)
	rand.Read(buf)

	// Хэшируем строку
	hash := sha256.New()
	hash.Write(buf)

	// refreshTokenString := strings.ToUpper(strings.TrimRight(base64.URLEncoding.EncodeToString(hash.Sum(nil)), "="))
	// Конвертируем в hex строку
	refreshTokenString := strings.ToLower(fmt.Sprintf("%x", hash.Sum(nil)))

	// Формируем новый refresh token
	refreshToken := models.RefreshToken{
		Token:   refreshTokenString,
		Expires: time.Now().Add(time.Hour * 24 * 10).Unix(), // 10 дней
	}

	return &refreshToken
}

// NewRefreshToken генерирует новый refresh token
func NewRefreshToken(ctx context.Context) (*models.RefreshToken, error) {
	// Генерируем новый refresh token
	refreshToken := generateRefreshToken()

	// Извлекаем из контекста пользователя
	user, ok := models.FromContextWithUser(ctx)
	if ok == false {
		log.Println("User is not found in context")
		return nil, errors.New("User is not found in context")
	}

	// Прописываем userID, которому будет принадлежать токен
	refreshToken.UserID = user.ID

	// Извлекаем из контекста refreshTokenStore
	refreshTokenStore, ok := store.FromContextWithRefreshTokenStore(ctx)
	if ok == false {
		log.Println("Refresh store is not found in context")
		return nil, errors.New("Refresh store is not found in context")
	}

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
		log.Printf("Error finding refresh token: %v\n", err)
		return nil, errors.New("Error finding refresh token")
	}
	now := time.Now().Unix()
	if now >= oldRefreshToken.Expires {
		// Токен истек
		log.Printf("The token has expired: now: %d, token_exp: %d\n", now, oldRefreshToken.Expires)
		return nil, errors.New("The token has expired")
	}

	// Генерируем новый refresh token
	refresh := generateRefreshToken()

	// Обновляем старый, активный токен
	updatedRefresh, err := refreshTokenStore.Update(token, refresh)
	if err != nil {
		log.Printf("Error updating refresh token: %v\n", err)
		return nil, errors.New("Error updating refresh token")
	}

	return updatedRefresh, nil
}
