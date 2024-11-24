package repository

import (
	"context"
	"time"
)

// Интерфейс репозитория
type ChatRepository interface {
	CreateChat(ctx context.Context, usernames []string) (int64, error)
	DeleteChat(ctx context.Context, chatID int64) error
	AddMessage(ctx context.Context, from string, text string, timestamp time.Time) error
	ChatExists(ctx context.Context, chatID int64) (bool, error) // Новый метод
}
