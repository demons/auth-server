package tokgen

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"strings"
)

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
