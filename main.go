package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"

	"audiolang.com/auth-server/notify"
	"audiolang.com/auth-server/senders"
	"audiolang.com/auth-server/store"
	"audiolang.com/auth-server/tokgen"

	// Driver postgres
	_ "github.com/lib/pq"
)

var (
	database *sql.DB

	// Доступ к хранилищу
	userDb         store.UserStore
	refreshTokenDb store.TokenStore
	tempTokenDb    store.TokenStore

	// Генерация jwt токенов
	jwtGen             *tokgen.JwtAccessGenerate
	jwtConfig          tokgen.Config
	tokenGenerator     *tokgen.TokenGenerator
	tempTokenGenerator *tokgen.TokenGenerator

	// Рассылка уведомлений
	emailConfig             senders.EmailConfig
	emailSender             *senders.EmailSender
	accountActivateTemplate *template.Template
	emailNotificator        *notify.EmailNotificator
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

	// Создадим таблицу временных токенов, если она отсутствует
	query, err = ioutil.ReadFile("./sql/create_tokens_table.sql")
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

	jwtConfig = tokgen.Config{
		Expires:    3600,
		PrivateKey: pKey,
	}

	emailConfig = senders.EmailConfig{
		Host:     "email-smtp.eu-west-1.amazonaws.com",
		Port:     587,
		Login:    "AKIAIFNXYQAACAXA3XTQ",
		Password: "AkRc59lK6KQztv8Y9710ereFE/tmu0XTaT1Sz1mVJ6rh",
	}

	// Загрузим template содержимого активации аккаунта
	accountActivateTemplate, err = template.ParseFiles("./notify/templates/account_activation.html")
	if err != nil {
		log.Fatalf("Error reading notify template: %v", err)
	}
}

func main() {

	userDb = store.NewUserDb(database)
	refreshTokenDb = store.NewRefreshTokenDb(database)
	tempTokenDb = store.NewTempTokenDb(database)

	jwtGen = jwtConfig.New()
	tokenGenerator = tokgen.NewTokenGenerator(time.Duration(time.Hour * 24 * 10))
	tempTokenGenerator = tokgen.NewTokenGenerator(time.Duration(time.Hour * 24))

	emailSender = emailConfig.NewEmailSender("notify@audiolang.com")
	emailNotificator = notify.NewEmailNotificator(emailSender)

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
