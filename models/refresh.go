package models

import (
	"time"
)

// RefreshToken содерржит информацию о refresh токене
type RefreshToken struct {
	UserID    int64
	Token     string
	Expires   int64
	CreatedAt time.Time
	UpdatedAt time.Time
}
