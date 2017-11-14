package store

import (
	"database/sql"

	"audiolang.com/auth-server/models"
)

// TempTokenDb хранилище для работы с пользователями
type TempTokenDb struct {
	db *sql.DB
}

// NewTempTokenDb конструктор
func NewTempTokenDb(database *sql.DB) *TempTokenDb {
	return &TempTokenDb{
		db: database,
	}
}

// Insert добавляет новый temp token в базу данных
func (t *TempTokenDb) Insert(token *models.Token) error {
	_, err := t.db.Exec("INSERT INTO tokens (user_id, token, expires, scopes) VALUES ($1, $2, $3, $4)", token.UserID, token.Token, token.Expires, token.Scopes)
	if err != nil {
		return err
	}

	return nil
}

// Update обновляет refresh token с указанным token
func (t *TempTokenDb) Update(token string, newToken *models.Token) error {
	_, err := t.db.Exec("UPDATE tokens SET token=$1, expires=$2, scopes=$3, updated_at=(now() at time zone 'utc') WHERE token=$4", newToken.Token, newToken.Expires, newToken.Scopes, token)
	if err != nil {
		return err
	}

	return nil
}

// FindByField выполняет поиск по указанному полю
func (t *TempTokenDb) FindByField(field string, value interface{}) (*models.Token, error) {
	var token models.Token

	err := t.db.QueryRow("SELECT user_id, token, expires, scopes FROM tokens WHERE "+field+"=$1", value).Scan(&token.UserID, &token.Token, &token.Expires, &token.Scopes)

	switch {
	case err == sql.ErrNoRows:
		// Элемент не найден
		return nil, nil
	case err != nil:
		// Произошла какая-то ошибка
		return nil, err
	}

	return &token, nil
}

// FindByUserID returns a refresh by userID
func (t *TempTokenDb) FindByUserID(userID int64) (*models.Token, error) {
	return t.FindByField("user_id", userID)
}

// FindByToken returns a refresh by RefreshToken
func (t *TempTokenDb) FindByToken(token string) (*models.Token, error) {
	return t.FindByField("token", token)
}
