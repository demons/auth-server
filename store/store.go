package store

import (
	"context"

	"audiolang.com/auth-server/models"
)

// The key type is unexported to prevent collisions with context keys defined in
// other packages.
type key int

// userIPkey is the context key for the user IP address.  Its value of zero is
// arbitrary.  If this package defined other context keys, they would have
// different integer values.
const (
	userStoreKey  key = 0
	tokenStoreKey key = 1
)

// UserStore для работы с хранилищем пользователей
type UserStore interface {
	Insert(*models.User) (int64, error)
	UpdatePassword(*models.User) error
	FindByField(string, interface{}) (*models.User, error)
	FindByEmail(string) (*models.User, error)
	FindByUserID(int64) (*models.User, error)
	FindBySID(string) (*models.User, error)
}

// TokenStore для работы с хранилищем токенов
type TokenStore interface {
	Insert(*models.Token) error
	Update(string, *models.Token) error
	FindByField(string, interface{}) (*models.Token, error)
	FindByUserID(int64) (*models.Token, error)
	FindByToken(string) (*models.Token, error)
}

// NewContextWithUserStore returns a new Context carrying user store.
func NewContextWithUserStore(ctx context.Context, store UserStore) context.Context {
	return context.WithValue(ctx, userStoreKey, store)
}

// FromContextWithUserStore extracts the user context from ctx, if present.
func FromContextWithUserStore(ctx context.Context) (UserStore, bool) {
	// ctx.Value returns nil if ctx has no value for the key;
	// the User type assertion returns ok=false for nil.
	store, ok := ctx.Value(userStoreKey).(UserStore)
	return store, ok
}

// NewContextWithTokenStore returns a new Context carrying token store.
func NewContextWithTokenStore(ctx context.Context, store TokenStore) context.Context {
	return context.WithValue(ctx, tokenStoreKey, store)
}

// FromContextWithTokenStore extracts the token store context from ctx, if present.
func FromContextWithTokenStore(ctx context.Context) (TokenStore, bool) {
	// ctx.Value returns nil if ctx has no value for the key;
	// the RefreshToken type assertion returns ok=false for nil.
	store, ok := ctx.Value(tokenStoreKey).(TokenStore)
	return store, ok
}
