package repository

import (
	_ "awesomeProject3/models"
	"database/sql"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateOrIgnore(telegramID int64, username string) error {
	_, err := r.db.Exec(
		"INSERT INTO users (telegram_id, username) VALUES ($1, $2) ON CONFLICT (telegram_id) DO NOTHING",
		telegramID, username,
	)
	return err
}

func (r *UserRepository) GetIDByTelegramID(telegramID int64) (int, error) {
	var id int
	err := r.db.QueryRow("SELECT id FROM users WHERE telegram_id = $1", telegramID).Scan(&id)
	return id, err
}
