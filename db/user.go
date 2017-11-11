package db

import (
	"context"
	"database/sql"

	"audiolang.com/auth-server/models"
)

// User хранилище для работы с пользователями
type User struct {
	db *sql.DB
}

// NewUserStore конструктор
func NewUserStore(database *sql.DB) *User {
	return &User{
		db: database,
	}
}

// Insert добавляет нового пользователя
func (u *User) Insert(user *models.User) (int64, error) {
	var userID int64
	err := u.db.QueryRow("INSERT INTO users (email, hash, salt, sid, name, is_social, is_active, is_verified) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id", user.Email, user.Hash, user.Salt, user.SID, user.Name, user.IsSocial, user.IsActive, user.IsVerified).Scan(&userID)
	if err != nil {
		return 0, err
	}

	return userID, nil
}

// FindByField выполняет поиск по указанному полю
func (u *User) FindByField(field string, value interface{}) (*models.User, error) {
	var user models.User
	err := u.db.QueryRow("SELECT id, email, hash, salt, sid, name, is_social, is_active, is_verified FROM users WHERE "+field+"=$1", value).Scan(&user.ID, &user.Email, &user.Hash, &user.Salt, &user.SID, &user.Name, &user.IsSocial, &user.IsActive, &user.IsVerified)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// FindByEmail returns a user by email
func (u *User) FindByEmail(email string) (*models.User, error) {
	return u.FindByField("email", email)
}

// FindByUserID returns a user by userID
func (u *User) FindByUserID(userID int64) (*models.User, error) {
	return u.FindByField("ID", userID)
}

// FindBySID return a user by social ID
func (u *User) FindBySID(sid string) (*models.User, error) {
	return u.FindByField("sid", sid)
}

// NewContextWithUserStore returns a new Context carrying user store.
func NewContextWithUserStore(ctx context.Context, user *User) context.Context {
	return context.WithValue(ctx, userStoreKey, user)
}

// FromContextWithUserStore extracts the user context from ctx, if present.
func FromContextWithUserStore(ctx context.Context) (*User, bool) {
	// ctx.Value returns nil if ctx has no value for the key;
	// the User type assertion returns ok=false for nil.
	user, ok := ctx.Value(userStoreKey).(*User)
	return user, ok
}
