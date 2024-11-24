package repository

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"time"
)

type PostgresChatRepository struct {
	db *pgx.Conn
}

// Конструктор
func NewPostgresChatRepository(db *pgx.Conn) *PostgresChatRepository {
	return &PostgresChatRepository{db: db}
}

// Реализация метода ChatExists
func (r *PostgresChatRepository) ChatExists(ctx context.Context, chatID int64) (bool, error) {
	var exists bool
	err := r.db.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM chats WHERE id = $1)", chatID).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}
func (r *PostgresChatRepository) CreateChat(ctx context.Context, usernames []string) (int64, error) {
	var chatID int64

	// Создаем новый чат
	err := r.db.QueryRow(ctx, "INSERT INTO chats DEFAULT VALUES RETURNING id").Scan(&chatID)
	if err != nil {
		return 0, fmt.Errorf("failed to create chat: %v", err)
	}

	// Добавляем пользователей к чату
	for _, username := range usernames {
		_, err = r.db.Exec(ctx, "INSERT INTO chat_users (chat_id, user_name) VALUES ($1, $2)", chatID, username)
		if err != nil {
			return 0, fmt.Errorf("failed to add user %s to chat: %v", username, err)
		}
	}

	return chatID, nil
}

func (r *PostgresChatRepository) DeleteChat(ctx context.Context, chatID int64) error {
	res, err := r.db.Exec(ctx, "DELETE FROM chats WHERE id = $1", chatID)
	if err != nil {
		return fmt.Errorf("failed to delete chat: %v", err)
	}
	if res.RowsAffected() == 0 {
		return fmt.Errorf("chat with id %d not found", chatID)
	}
	return nil
}

func (r *PostgresChatRepository) AddMessage(ctx context.Context, from string, text string, timestamp time.Time) error {
	_, err := r.db.Exec(ctx, "INSERT INTO message (from_user, message, timestamp) VALUES ($1, $2, $3)", from, text, timestamp)
	if err != nil {
		return fmt.Errorf("failed to send message: %v", err)
	}
	return nil
}
