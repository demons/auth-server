package utils

import (
	"encoding/base64"
	"errors"
	"strings"

	"golang.org/x/crypto/scrypt"
)

// HashPassword создает хэш пароля
func HashPassword(password string, salt []byte) (string, error) {
	hash, err := scrypt.Key([]byte(password), salt, 1<<14, 8, 1, 64)
	if err != nil {
		return "", err
	}

	return strings.TrimRight(base64.StdEncoding.EncodeToString(hash), "="), nil
}

// VerifyPassword проверяем корректный ли пароль
func VerifyPassword(hash, salt, password string) error {
	if hash == "" || salt == "" {
		return errors.New("Incorrect password")
	}
	// Декодируем соль в массив байтов
	s, err := base64.StdEncoding.DecodeString(salt + "==")
	if err != nil {
		return errors.New("Bad salt")
	}

	// Хэшируем пароль
	verifyHash, err := HashPassword(password, s)
	if err != nil {
		return err
	}

	// Проверяем совпадают ли хэши
	if hash == verifyHash {
		return nil
	}

	return errors.New("Incorrect password")
}
