package oauth

import "errors"

// Provider interface
type Provider interface {
	ExchangeCode(code string) (string, error)
	GetUserProfile(accessToken string) (*UserProfile, error)
}

var providerMap = map[string]Provider{
	"vk": VkProvider{},
	// "facebook": FacebookProvider{},
	// "google": GoogleProvider{},
}

// GetProviderByName returns provider by name
func GetProviderByName(providerName string) (Provider, error) {
	provider := providerMap[providerName]
	if provider == nil {
		return nil, errors.New("This provider is not supported")
	}

	return provider, nil
}

// TODO: Реализовать доступ к сервисам соц. сетей

// Меняем code на access token
// Получаем профиль пользователя из соц. сети
// Отдаем результат
