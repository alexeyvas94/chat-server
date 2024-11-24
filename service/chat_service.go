package service

import (
	"context"
	"time"
)

// Интерфейс ChatService
type ChatService interface {
	CreateChat(ctx context.Context, usernames []string) (int64, error)
	DeleteChat(ctx context.Context, chatID int64) error
	AddMessage(ctx context.Context, from, text string, timestamp time.Time) error
}
