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
	Scopes  string
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

// SetScope добавляет права токену
func (t *Token) SetScope(scope string) {
	if t.Scopes == "" {
		t.Scopes = scope
		return
	}
	strings.TrimRight(t.Scopes, ";")
	t.Scopes = t.Scopes + ";" + scope
}

// SetScopes добавляет несколько прав токену
func (t *Token) SetScopes(scopes []string) {
	for _, v := range scopes {
		t.SetScope(v)
	}
}

// GetScope проверяет имеется ли указанный scope
func (t *Token) GetScope(scope string) bool {
	scopes := strings.Split(t.Scopes, ";")
	for _, v := range scopes {
		if v == scope {
			return true
		}
	}

	return false
}

// Valid проверяет валидность токена
func (t *Token) Valid() bool {
	now := time.Now().Unix()
	if now >= t.Expires {
		// Токен истек
		return false
	}

	return true
}
