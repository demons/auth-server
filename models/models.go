package models

import "context"

// The key type is unexported to prevent collisions with context keys defined in
// other packages.
type key int

// userIPkey is the context key for the user IP address.  Its value of zero is
// arbitrary.  If this package defined other context keys, they would have
// different integer values.
const (
	userKey key = 0
)

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
