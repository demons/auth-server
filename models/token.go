package models

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"strings"
	"time"
)

// Token модель для хранения токена
type Token struct {
	UserID  int64
	Token   string
	Expires int64
}

// generateToken генерирует случайный токен
func generateToken() string {
	// Последовательность случайных байт
	buf := make([]byte, 128)
	rand.Read(buf)

	// Хэшируем строку
	hash := sha256.New()
	hash.Write(buf)

	// Конвертируем в hex строку
	tokenString := strings.ToLower(fmt.Sprintf("%x", hash.Sum(nil)))

	return tokenString
}

// NewToken создает новый токен
func NewToken(userID int64, expireIn time.Duration) *Token {
	return &Token{
		UserID:  userID,
		Token:   generateToken(),
		Expires: time.Now().Add(expireIn).Unix(),
	}
}

// Valid проверяет валидность токена
func (t Token) Valid() bool {
	now := time.Now().Unix()
	if now >= t.Expires {
		// Токен истек
		return false
	}

	return true
}
