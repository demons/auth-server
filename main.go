package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
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
	configPath string
	config     *Config

	database *sql.DB

	// Доступ к хранилищу
	userDb         store.UserStore
	refreshTokenDb store.TokenStore
	tempTokenDb    store.TokenStore

	// Генерация jwt токенов
	jwtGen             *tokgen.JwtAccessGenerate
	jwtConfig          tokgen.Config
	tokenGenerator     *tokgen.TokenGenerator
	tokenChecker       *tokgen.JwtAccessChecker
	tempTokenGenerator *tokgen.TokenGenerator

	// Рассылка уведомлений
	emailConfig      senders.EmailConfig
	emailSender      *senders.EmailSender
	emailNotificator *notify.EmailNotificator
)

// Шаблоны для электронных писем
var messageTemplatePaths = map[string]string{
	"activateAccount": "./notify/templates/activate_account.html",
	"resetPassword":   "./notify/templates/reset_password.html",
}
var messageTemplates = make(map[string]*template.Template, 3)

func init() {

	flag.StringVar(&configPath, "config", "config.json", "path to config file")

	config = loadConfiguration(configPath)

	databaseInit()
	tokensInit()

	emailConfig = senders.EmailConfig{
		Host:     config.Email.Host,
		Port:     config.Email.Port,
		Login:    config.Email.Login,
		Password: config.Email.Password,
	}

	// Загрузим шаблоны для текстов писем
	for k, v := range messageTemplatePaths {
		parsedTemplate, err := template.ParseFiles(v)
		if err != nil {
			log.Fatalf("Error reading notify template: %v", err)
		}
		messageTemplates[k] = parsedTemplate
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

	// PUBLIC APIs
	router.POST("/token", HandleToken)

	// Регистрация нового пользователя
	router.POST("/account/signup", HandleSignUp)

	// Сброс пароля
	router.POST("/account/password/reset", HandlePasswordReset)

	// PRIVATE APIs

	// TODO: Закрытые APIs должены быть доступены только аутентифицированным пользователям

	// Смена пароля
	router.POST("/account/password/change", Auth(HandlePasswordChange))

	// Запускаем сервер
	fmt.Println("Server started...", fmt.Sprintf("%s:%d", config.Host, config.Port))
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%d", config.Host, config.Port), router))
}

// Config for configuration application
type Config struct {
	Host  string `json:"host"`
	Port  int    `json:"port"`
	Email struct {
		Host     string `json:"host"`
		Port     int    `json:"port"`
		Login    string `json:"login"`
		Password string `json:"password"`
	} `json:"email"`
	Database struct {
		Host     string `json:"host"`
		Port     int    `json:"port"`
		Login    string `json:"login"`
		Password string `json:"password"`
		DbName   string `json:"db"`
	} `json:"database"`
}

func loadConfiguration(file string) *Config {
	var config Config
	configFile, err := os.Open(file)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer configFile.Close()
	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&config)
	return &config
}

func databaseInit() {
	var err error
	// Инициализируем базу данных
	database, err = sql.Open("postgres", fmt.Sprintf("host='%s' port=%d user=%s password=%s dbname=%s sslmode=disable", config.Database.Host, config.Database.Port, config.Database.Login, config.Database.Password, config.Database.DbName))
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
}

func tokensInit() {
	// Считываем приватный ключ
	pKey, err := ioutil.ReadFile("./secrets/app.rsa")
	if err != nil {
		log.Fatalf("Error reading private key: %v\n", err)
	}

	pubKey, err := ioutil.ReadFile("./secrets/app.rsa.pub")
	if err != nil {
		log.Fatalf("Error reading public key: %v", err)
	}

	jwtConfig = tokgen.Config{
		Expires:    3600,
		PrivateKey: pKey,
	}

	tokenChecker = tokgen.NewJwtAccessChecker(pubKey)
}
