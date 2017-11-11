package oauth

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/vk"
)

// TODO: Сделать загрузку этих данных из общего конфига
var config = &oauth2.Config{
	RedirectURL:  "http://localhost:3000/vk-auth",
	ClientID:     "6163381",
	ClientSecret: "jeBgF2a9ynznseZPfeqz",
	Scopes:       []string{""},
	Endpoint:     vk.Endpoint,
}

// VkProvider provider
type VkProvider struct {
}

type user struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type responseJSON struct {
	Users []user `json:"response"`
}

// ExchangeCode обменивает code на access_token
func (p VkProvider) ExchangeCode(code string) (string, error) {
	// Обмениваем code на token
	token, err := config.Exchange(oauth2.NoContext, code)
	if err != nil {
		log.Printf("Code exchange failed: %v\n", err)
		return "", errors.New("Code exchange failed")
	}

	return token.AccessToken, nil
}

// GetUserProfile returns user profile by access token
func (p VkProvider) GetUserProfile(accessToken string) (*UserProfile, error) {
	// Делаем запрос на получение информации о пользователе
	response, err := http.Get("https://api.vk.com/method/users.get?v=5.68&access_token=" + accessToken)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	// Обрабатываем полученную информацию о пользователе
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var resp responseJSON

	err = json.Unmarshal(body, &resp)
	if err != nil {
		return nil, err
	}

	if len(resp.Users) == 0 {
		return nil, errors.New("Reply does not contain a user")
	}

	name := resp.Users[0].FirstName + " " + resp.Users[0].LastName

	return &UserProfile{
		ID:           strconv.Itoa(resp.Users[0].ID),
		Name:         name,
		ProviderName: "vk",
	}, nil
}
