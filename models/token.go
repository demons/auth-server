package models

import (
	"time"
)

// Token содержит информацию о паре access и refresh токенах
type Token struct {
	Access           string        `json:"access_token"`
	Refresh          string        `json:"refresh_token"`
	RefreshExpiresIn time.Duration `json:"refresh_exp"`
}

// RefreshToken содерржит информацию о refresh токене
type RefreshToken struct {
	UserID    int64
	Token     string
	Expires   int64
	CreatedAt time.Time
	UpdatedAt time.Time
}
