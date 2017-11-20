package tokgen

import (
	"errors"
	"log"

	jwt "github.com/dgrijalva/jwt-go"
)

var (
	// ErrTokenIsInvalid is error const
	ErrTokenIsInvalid = errors.New("Token is not valid")
)

// JwtAccessChecker валидатор JWT токена
type JwtAccessChecker struct {
	PublicKey []byte
}

// NewJwtAccessChecker конструктор
func NewJwtAccessChecker(publicKey []byte) *JwtAccessChecker {
	return &JwtAccessChecker{
		PublicKey: publicKey,
	}
}

// Check is checking token string
func (ac *JwtAccessChecker) Check(tokenString string) (map[string]interface{}, error) {
	pubKey, err := jwt.ParseRSAPublicKeyFromPEM(ac.PublicKey)
	if err != nil {
		log.Printf("Error parsing public key: %v", err)
		return nil, err
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return pubKey, nil
	})

	if err != nil {
		log.Printf("Error parsing token: %v", err)
		return nil, err
	}

	if token.Valid == true {
		return token.Claims.(jwt.MapClaims), nil
	}

	return nil, ErrTokenIsInvalid
}
