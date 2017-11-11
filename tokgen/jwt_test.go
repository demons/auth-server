package tokgen

import (
	"testing"
)

// TODO: Реализовать тесты генератора токена доступа

// Путь к закрытому ключу
const (
	privateKeyPath string = "../app.rsa"
	publicKeyPath  string = "../app.rsa.pub"
	userID         int64  = 25
)

func SetUp(t *testing.T) {
	// // Читаем из файла private key
	// pKey, err := ioutil.ReadFile(privateKeyPath)
	// if err != nil {
	// 	t.Fatal("Ошибка при загрузке закрытого ключа")
	// }

	// // Преобразовываем закрытый ключ
	// privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(pKey)
	// if err != nil {
	// 	t.Error("Ошибка при преобразовании закрытого ключа")
	// }
	t.Log("Init")
}

func TestGetToken(t *testing.T) {

	// // create a rsa 256 signer
	// signer := jwt.New(jwt.GetSigningMethod("RS256"))

	// // set claims
	// claims := signer.Claims.(jwt.MapClaims)
	// claims["uid"] = userID
	// claims["exp"] = time.Date(2017, 10, 29, 20, 0, 0, 0, time.UTC)

	// // Подписываем токен закрытым ключем
	// tokenString, err := signer.SignedString(privateKey)
	// if err != nil {
	// 	t.Error("Ошибка при подписывании токена")
	// }

	// // Создаем класс JwtToken
	// jgen, err := NewJwtToken(pKey)
	// if err != nil {
	// 	t.Error("Ошибка при создании JwtGenerator")
	// }
	// jgen.Payload["uid"] = userID
	// token, err := jgen.GetToken()
	// if err != nil {
	// 	t.Error("Ошибка при генерации jwt токена")
	// }

	// if tokenString != token {
	// 	t.Error("Jwt токены не равны")
	// }
}

// func TestCheckToken(t *testing.T) {
// 	// Путь к открытому ключу
// 	publicKeyPath := "../app.rsa.pub"

// 	// Читаем из файла public key
// 	pubKey, err := ioutil.ReadFile(publicKeyPath)
// 	if err != nil {
// 		t.Error("Ошибка при загрузке открытого ключа")
// 	}
// 	jtoken, err := JwtToken(pKey, pubKey)
// }

func TestCheckRefreshToken(t *testing.T) {

}
