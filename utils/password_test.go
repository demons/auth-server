package utils

import (
	"encoding/base64"
	"strings"
	"testing"

	"golang.org/x/crypto/scrypt"
)

func TestHashPassword(t *testing.T) {
	// Создаем соль
	salt := []byte{1, 2, 3, 4}

	// Пароль
	password := "testPassword"

	h, _ := scrypt.Key([]byte(password), salt, 1<<14, 8, 1, 64)
	expected := strings.Trim(base64.StdEncoding.EncodeToString(h), "=")

	hash, _ := HashPassword(password, salt)

	if expected != hash {
		t.Error("Хэши не совпадают")
	}

}

func TestVerifyPassword(t *testing.T) {
	// VerifyPassword(hash, salt, password) error
	// Хэшь пароля `testPassword`
	passwordHash := "OY9dS/J0KkNOSZE3Fw+Yg5RkRScIOZPISXhgKjo1eOgeDuSRuwaeUBQ84YKW2YHlB1rnXC1Qd5mS/NFSer/wbA"
	salt := []byte{1, 2, 3, 4}
	saltString := strings.TrimRight(base64.StdEncoding.EncodeToString(salt), "=")

	truePassword := "testPassword"
	// falsePassword := "testPasswordFalse"

	// Хэши паролей должни совпадать
	if err := VerifyPassword(passwordHash, saltString, truePassword); err != nil {
		t.Error("Неправильный пароль, ошибок быть не должно")
	}

	// Ожидаем ошибку
	if err := VerifyPassword(passwordHash, saltString, "falsePassword1234"); err == nil {
		t.Errorf("Должна быть ошибка некорректного пароля: %v", err)
	}

	// Хэш и соль пустые, ждем ошибку
	if err := VerifyPassword("", "", truePassword); err == nil {
		t.Errorf("Должна быть ошибка некорректного пароля, т.к. хэш и соль пустые: %v", err)
	}

	if err := VerifyPassword("", saltString, truePassword); err == nil {
		t.Errorf("Должна быть ошибка некорректного пароля, т.к. хэш пустой: %v", err)
	}

	if err := VerifyPassword(passwordHash, "", truePassword); err == nil {
		t.Errorf("Должна быть ошибка некорректного пароля, т.к. соль пустая: %v", err)
	}
}
