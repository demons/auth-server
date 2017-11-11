package models

import (
	"context"
	"database/sql"
)

// User содержит информацию о пользователе
type User struct {
	ID         int64
	Email      sql.NullString
	Hash       sql.NullString
	Salt       sql.NullString
	IsVerified sql.NullBool
	IsActive   bool
	IsSocial   bool
	SID        sql.NullString
	Name       sql.NullString
}

// NewContextWithUser returns a new Context carrying user.
func NewContextWithUser(ctx context.Context, user *User) context.Context {
	return context.WithValue(ctx, userKey, user)
}

// FromContextWithUser extracts the user from ctx, if present.
func FromContextWithUser(ctx context.Context) (*User, bool) {
	// ctx.Value returns nil if ctx has no value for the key;
	// the User type assertion returns ok=false for nil.
	user, ok := ctx.Value(userKey).(*User)
	return user, ok
}
