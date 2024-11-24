package service

import (
	"context"
	"errors"
	"time"

	"github.com/alexeyvas94/chat-server/repository"
)

type chatService struct {
	chatRepo repository.ChatRepository
}

func NewChatService(chatRepo repository.ChatRepository) ChatService {
	return &chatService{chatRepo: chatRepo}
}

func (s *chatService) CreateChat(ctx context.Context, usernames []string) (int64, error) {
	// проверяем, что список пользователей не пуст
	if len(usernames) == 0 {
		return 0, ErrEmptyUserList
	}
	return s.chatRepo.CreateChat(ctx, usernames)
}

func (s *chatService) DeleteChat(ctx context.Context, chatID int64) error {
	// проверка на существование чата
	exists, err := s.chatRepo.ChatExists(ctx, chatID)
	if err != nil {
		return err
	}
	if !exists {
		return ErrChatNotFound
	}
	return s.chatRepo.DeleteChat(ctx, chatID)
}

func (s *chatService) AddMessage(ctx context.Context, from, text string, timestamp time.Time) error {
	// валидация данных сообщения
	if from == "" || text == "" {
		return ErrInvalidMessage
	}
	return s.chatRepo.AddMessage(ctx, from, text, timestamp)
}

// Определите ошибки
var (
	ErrEmptyUserList  = errors.New("user list is empty")
	ErrChatNotFound   = errors.New("chat not found")
	ErrInvalidMessage = errors.New("invalid message data")
)
