package tokgen

import (
	"context"
	"errors"
	"log"
	"time"

	"audiolang.com/auth-server/models"
	"audiolang.com/auth-server/store"
)

// The key type is unexported to prevent collisions with context keys defined in
// other packages.
type key int

// userIPkey is the context key for the user IP address.  Its value of zero is
// arbitrary.  If this package defined other context keys, they would have
// different integer values.
const (
	tokenKey key = 0
)

// TokenGenerator занимается созданием и обновлением токенов
type TokenGenerator struct {
	expireIn time.Duration
}

// NewTokenGenerator создает новый генератор токенов
func NewTokenGenerator(expireIn time.Duration) *TokenGenerator {
	return &TokenGenerator{
		expireIn: expireIn,
	}
}

// FindToken возвращает токен из хранилища
func (g TokenGenerator) FindToken(ctx context.Context, token string) (*models.Token, error) {
	// Получить хранилище токенов из ctx
	tokenStore, ok := store.FromContextWithTokenStore(ctx)
	if ok == false {
		log.Println("Token store is not found in context")
		return nil, errors.New("Token store is not found in context")
	}

	// Выполнить поиск токена в хранилище
	findedToken, err := tokenStore.FindByToken(token)
	if err != nil {
		// Произошла какая-то ошибка при поиске токена
		log.Printf("Error finding token: %v\n", err)
		return nil, errors.New("Error finding token")
	}

	if findedToken == nil {
		return nil, nil
	}

	return findedToken, nil
}

// CreateToken создает новый токен, если его еще не существует в хранилище
func (g TokenGenerator) CreateToken(ctx context.Context, scopes []string) (*models.Token, error) {
	// Получить хранилище токенов из ctx
	tokenStore, ok := store.FromContextWithTokenStore(ctx)
	if ok == false {
		log.Println("Token store is not found in context")
		return nil, errors.New("Token store is not found in context")
	}

	// Получить текущего пользователя из ctx
	user, ok := models.FromContextWithUser(ctx)
	if ok == false {
		log.Println("User is not found in context")
		return nil, errors.New("User is not found in context")
	}

	// Генерируем новый токен для этого пользователя
	newToken := models.NewToken(user.ID, g.expireIn)
	newToken.SetScopes(scopes)

	// Сохраняем токен в хранилище
	err := tokenStore.Insert(newToken)
	if err != nil {
		log.Printf("Error inserting token in the store: %v", err)
		return nil, errors.New("Error inserting token in the store")
	}

	return newToken, nil
}

// UpdateToken обновляет токен
func (g TokenGenerator) UpdateToken(ctx context.Context, token *models.Token) (*models.Token, error) {
	// Проверим валиден ли токен
	if token.Valid() == false {
		log.Println("This token is not valid")
		return nil, errors.New("This token is not valid")
	}

	// Получить хранилище токенов из ctx
	tokenStore, ok := store.FromContextWithTokenStore(ctx)
	if ok == false {
		log.Println("Token store is not found in context")
		return nil, errors.New("Token store is not found in context")
	}

	newToken := models.NewToken(token.UserID, g.expireIn)

	// Обновляем токен в хранилище
	err := tokenStore.Update(token.Token, newToken)
	if err != nil {
		log.Printf("Error updating token: %v\n", err)
		return nil, errors.New("Error updating token")
	}

	return newToken, nil
}

// Delete удаляет токен
// func (t Token) Delete(ctx context.Context) error {

// }

// NewContextWithTokenGenerator returns a new Context carrying token generator.
func NewContextWithTokenGenerator(ctx context.Context, tokenGenerator *TokenGenerator) context.Context {
	return context.WithValue(ctx, tokenKey, tokenGenerator)
}

// FromContextWithTokenGenerator extracts the token generator from ctx, if present.
func FromContextWithTokenGenerator(ctx context.Context) (*TokenGenerator, bool) {
	// ctx.Value returns nil if ctx has no value for the key;
	// the User type assertion returns ok=false for nil.
	tokenGenerator, ok := ctx.Value(tokenKey).(*TokenGenerator)
	return tokenGenerator, ok
}
