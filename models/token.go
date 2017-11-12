package models

import "time"

// Token модель для хранения токена
type Token struct {
	UserID    int64
	Token     string
	Expires   int64
	CreatedAt time.Time
	UpdatedAt time.Time
}
