package tokgen

import (
	"context"
	"errors"
	"log"
	"time"

	"audiolang.com/auth-server/models"
	jwt "github.com/dgrijalva/jwt-go"
)

// Config для генерации нового JWT токена
type Config struct {
	PrivateKey []byte
	// ID закрытого ключа, информация для расшифровки токена
	PrivateKeyID string
	Expires      time.Duration
}

// JwtAccessGenerate генератор JWT токена
type JwtAccessGenerate struct {
	cnf *Config
}

// New создает Token генератор
func (cf *Config) New() *JwtAccessGenerate {

	return &JwtAccessGenerate{
		cnf: cf,
	}
}

// Token генерирует новый токен
func (g JwtAccessGenerate) Token(ctx context.Context) (string, error) {
	// Парсим private key
	pKey, err := jwt.ParseRSAPrivateKeyFromPEM(g.cnf.PrivateKey)
	if err != nil {
		log.Printf("Error parsing private key: %v\n", err)
		return "", err
	}

	// create a rsa 256 signer
	signer := jwt.New(jwt.GetSigningMethod("RS256"))

	// set claims
	claims := signer.Claims.(jwt.MapClaims)

	// Записываем в payload идентификатор закрытого ключа
	if g.cnf.PrivateKeyID != "" {
		claims["PrivateKeyID"] = g.cnf.PrivateKeyID
	}

	// Получаем ID пользователя из ctx
	user, ok := models.FromContextWithUser(ctx)
	if ok == false {
		log.Printf("User is not found in context\n")
		return "", errors.New("User is not found in context")
	}
	// Записываем в payload идентификатор пользователя
	claims["userID"] = user.ID

	// Время действия access токена
	claims["exp"] = time.Now().Add(time.Second * g.cnf.Expires).UTC()

	// Подписываем jwt access token
	tokenString, err := signer.SignedString(pKey)
	if err != nil {
		log.Printf("Error signing token: %v\n", err)
		return "", err
	}

	return tokenString, nil
}
