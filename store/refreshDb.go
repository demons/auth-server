package store

import (
	"database/sql"

	"audiolang.com/auth-server/models"
)

// RefreshTokenDb хранилище для работы с пользователями
type RefreshTokenDb struct {
	db *sql.DB
}

// NewRefreshTokenDb конструктор
func NewRefreshTokenDb(database *sql.DB) *RefreshTokenDb {
	return &RefreshTokenDb{
		db: database,
	}
}

// Insert добавляет новый refresh token в базу данных
func (r *RefreshTokenDb) Insert(refToken *models.Token) error {
	_, err := r.db.Exec("INSERT INTO reftoks (user_id, token, expires) VALUES ($1, $2, $3)", refToken.UserID, refToken.Token, refToken.Expires)
	if err != nil {
		return err
	}

	return nil
}

// Update обновляет refresh token с указанным token
func (r *RefreshTokenDb) Update(token string, refToken *models.Token) error {

	_, err := r.db.Exec("UPDATE reftoks SET token=$1, expires=$2, updated_at=(now() at time zone 'utc') WHERE token=$3", refToken.Token, refToken.Expires, token)
	if err != nil {
		return err
	}

	return nil
}

// FindByField выполняет поиск по указанному полю
func (r *RefreshTokenDb) FindByField(field string, value interface{}) (*models.Token, error) {
	var refresh models.Token

	err := r.db.QueryRow("SELECT user_id, token, expires FROM reftoks WHERE "+field+"=$1", value).Scan(&refresh.UserID, &refresh.Token, &refresh.Expires)

	switch {
	case err == sql.ErrNoRows:
		// Элемент не найден
		return nil, nil
	case err != nil:
		// Произошла какая-то ошибка
		return nil, err
	}

	return &refresh, nil
}

// FindByUserID returns a refresh by userID
func (r *RefreshTokenDb) FindByUserID(userID int64) (*models.Token, error) {
	return r.FindByField("user_id", userID)
}

// FindByToken returns a refresh by RefreshToken
func (r *RefreshTokenDb) FindByToken(token string) (*models.Token, error) {
	return r.FindByField("token", token)
}
