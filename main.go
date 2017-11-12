package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"

	"audiolang.com/auth-server/store"
	"audiolang.com/auth-server/tokgen"

	// Driver postgres
	_ "github.com/lib/pq"
)

var (
	database       *sql.DB
	userDb         store.UserStore
	refreshTokenDb store.TokenStore
	jwtGen         *tokgen.JwtAccessGenerate
	jwtCnf         *tokgen.Config
	tokenGenerator *tokgen.TokenGenerator
)

func init() {
	var err error
	// Инициализируем базу данных
	database, err = sql.Open("postgres", "host='192.168.99.100' port=5432 user=postgres password=docker dbname=db sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	// Создадим таблицу пользователей, если она отсутствует
	query, err := ioutil.ReadFile("./sql/create_users_table.sql")
	if err != nil {
		log.Fatalf("Error reading sql: %v\n", err)
	}
	_, err = database.Exec(string(query))
	if err != nil {
		log.Fatalf("Error executing sql: %v\n", err)
	}

	// Создадим таблицу refresh токенов, если она отсутствует
	query, err = ioutil.ReadFile("./sql/create_reftoks_table.sql")
	if err != nil {
		log.Fatalf("Error reading sql: %v\n", err)
	}
	_, err = database.Exec(string(query))
	if err != nil {
		log.Fatalf("Error executing sql: %v\n", err)
	}

	// Считываем приватный ключ
	pKey, err := ioutil.ReadFile("./secrets/app.rsa")
	if err != nil {
		log.Fatalf("Error reading private key: %v\n", err)
	}

	jwtCnf = &tokgen.Config{
		Expires:    3600,
		PrivateKey: pKey,
	}
}

func main() {

	userDb = store.NewUserDb(database)
	refreshTokenDb = store.NewRefreshTokenDb(database)
	jwtGen = jwtCnf.New()
	tokenGenerator = tokgen.NewTokenGenerator(time.Duration(time.Hour * 24 * 10))

	router := httprouter.New()
	router.POST("/token", HandleToken)

	// Регистрация нового пользователя
	router.POST("/account/signup", HandleSignUp)

	// Сброс пароля
	// router.POST("/account/reset", HandleReset)

	// Смена пароля
	// router.POST("/account/password", HandleChangePassword)

	fmt.Println("Hello")
	log.Fatal(http.ListenAndServe("localhost:8000", router))
}
