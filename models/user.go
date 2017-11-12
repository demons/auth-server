package models

import (
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
