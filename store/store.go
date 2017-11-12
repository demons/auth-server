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
	userStoreKey    key = 0
	refreshStoreKey key = 1
)

// UserStore для работы с хранилищем пользователей
type UserStore interface {
	Insert(*models.User) (int64, error)
	FindByField(string, interface{}) (*models.User, error)
	FindByEmail(string) (*models.User, error)
	FindByUserID(int64) (*models.User, error)
	FindBySID(string) (*models.User, error)
}

// RefreshTokenStore для работы с хранилищем refresh токенов
type RefreshTokenStore interface {
	Insert(*models.RefreshToken) error
	Update(string, *models.RefreshToken) (*models.RefreshToken, error)
	FindByField(string, interface{}) (*models.RefreshToken, error)
	FindByUserID(int64) (*models.RefreshToken, error)
	FindByRefreshToken(string) (*models.RefreshToken, error)
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

// NewContextWithRefreshTokenStore returns a new Context carrying refresh store.
func NewContextWithRefreshTokenStore(ctx context.Context, store RefreshTokenStore) context.Context {
	return context.WithValue(ctx, refreshStoreKey, store)
}

// FromContextWithRefreshTokenStore extracts the refresh context from ctx, if present.
func FromContextWithRefreshTokenStore(ctx context.Context) (RefreshTokenStore, bool) {
	// ctx.Value returns nil if ctx has no value for the key;
	// the RefreshToken type assertion returns ok=false for nil.
	store, ok := ctx.Value(refreshStoreKey).(RefreshTokenStore)
	return store, ok
}
