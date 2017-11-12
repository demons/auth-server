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
func (r *RefreshTokenDb) Insert(refToken *models.RefreshToken) error {
	_, err := r.db.Exec("INSERT INTO reftoks (user_id, token, expires) VALUES ($1, $2, $3)", refToken.UserID, refToken.Token, refToken.Expires)
	if err != nil {
		return err
	}

	return nil
}

// Update обновляет refresh token с указанным token
func (r *RefreshTokenDb) Update(token string, refToken *models.RefreshToken) (*models.RefreshToken, error) {
	// Копируем старый токен, вернем пользователю новый токен
	updatedRefresh := *refToken

	err := r.db.QueryRow("UPDATE reftoks SET token=$1, expires=$2, updated_at=(now() at time zone 'utc') WHERE token=$3 RETURNING user_id, updated_at", refToken.Token, refToken.Expires, token).Scan(&updatedRefresh.UserID, &updatedRefresh.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &updatedRefresh, nil
}

// Deactivate деактивирует токен
// func (r *RefreshToken) Deactivate(token string) (int64, error) {
// 	var userID int64
// 	err := r.db.QueryRow("UPDATE reftoks SET is_active=$1 WHERE token=$2 AND is_active<>false RETURNING user_id", false, token).Scan(&userID)
// 	if err != nil {
// 		return 0, err
// 	}

// 	return userID, nil
// }

// FindByField выполняет поиск по указанному полю
func (r *RefreshTokenDb) FindByField(field string, value interface{}) (*models.RefreshToken, error) {
	var refresh models.RefreshToken
	err := r.db.QueryRow("SELECT user_id, token, expires, created_at, updated_at FROM reftoks WHERE "+field+"=$1", value).Scan(&refresh.UserID, &refresh.Token, &refresh.Expires, &refresh.CreatedAt, &refresh.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &refresh, nil
}

// FindByUserID returns a refresh by userID
func (r *RefreshTokenDb) FindByUserID(userID int64) (*models.RefreshToken, error) {
	return r.FindByField("user_id", userID)
}

// FindByRefreshToken returns a refresh by RefreshToken
func (r *RefreshTokenDb) FindByRefreshToken(token string) (*models.RefreshToken, error) {
	return r.FindByField("token", token)
}
