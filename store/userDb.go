package store

import (
	"database/sql"

	"audiolang.com/auth-server/models"
)

// UserDb хранилище для работы с пользователями
type UserDb struct {
	db *sql.DB
}

// NewUserDb конструктор
func NewUserDb(database *sql.DB) *UserDb {
	return &UserDb{
		db: database,
	}
}

// Insert добавляет нового пользователя
func (u *UserDb) Insert(user *models.User) (int64, error) {
	var userID int64
	err := u.db.QueryRow("INSERT INTO users (email, hash, salt, sid, name, is_social, is_active, is_verified) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id", user.Email, user.Hash, user.Salt, user.SID, user.Name, user.IsSocial, user.IsActive, user.IsVerified).Scan(&userID)
	if err != nil {
		return 0, err
	}

	return userID, nil
}

// UpdatePassword обновляет hash и salt
func (u *UserDb) UpdatePassword(user *models.User) error {
	_, err := u.db.Exec("UPDATE users SET hash=$1, salt=$2 WHERE id=$3", user.Hash, user.Salt, user.ID)
	if err != nil {
		return err
	}

	return nil
}

// FindByField выполняет поиск по указанному полю
func (u *UserDb) FindByField(field string, value interface{}) (*models.User, error) {
	var user models.User

	err := u.db.QueryRow("SELECT id, email, hash, salt, sid, name, is_social, is_active, is_verified FROM users WHERE "+field+"=$1", value).Scan(&user.ID, &user.Email, &user.Hash, &user.Salt, &user.SID, &user.Name, &user.IsSocial, &user.IsActive, &user.IsVerified)

	switch {
	case err == sql.ErrNoRows:
		// Элемент не найден
		return nil, nil
	case err != nil:
		// Произошла какая-то ошибка
		return nil, err
	}

	return &user, nil
}

// FindByEmail returns a user by email
func (u *UserDb) FindByEmail(email string) (*models.User, error) {
	return u.FindByField("email", email)
}

// FindByUserID returns a user by userID
func (u *UserDb) FindByUserID(userID int64) (*models.User, error) {
	return u.FindByField("ID", userID)
}

// FindBySID return a user by social ID
func (u *UserDb) FindBySID(sid string) (*models.User, error) {
	return u.FindByField("sid", sid)
}
