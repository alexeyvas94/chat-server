package main

import (
	"context"
	"fmt"
	pb "github.com/alexeyvas94/chat-server/proto"
	"github.com/jackc/pgx/v4"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
	"net"
	"sync"
)

const (
	dbDSN = "host=localhost port=54322 dbname=chat user=chat-user password=chat-password sslmode=disable"
)

type ChatServer struct {
	pb.UnimplementedChatServer
	mu sync.Mutex
}

func (c *ChatServer) Create(ctx context.Context, req *pb.CreateRequest) (*pb.CreateResponse, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	// Создаем соединение с базой данных
	con, err := pgx.Connect(ctx, dbDSN)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer con.Close(ctx)

	// Переменная для хранения id нового чата
	var chatID int64

	// Создаем новый чат и получаем его id
	err = con.QueryRow(ctx, "INSERT INTO chats DEFAULT VALUES RETURNING id").Scan(&chatID)
	if err != nil {
		return nil, fmt.Errorf("failed to create chat: %v", err)
	}
	// Вставляем пользователей в таблицу chat_users
	for _, userID := range req.Usernames {
		_, err = con.Exec(ctx, "INSERT INTO chat_users (chat_id, user_name) VALUES ($1, $2)", chatID, userID)
		if err != nil {
			return nil, fmt.Errorf("failed to add user %d to chat: %v", userID, err)
		}
	}
	// Возвращаем ответ с id созданного чата
	return &pb.CreateResponse{Id: chatID}, nil
}
func (c *ChatServer) Delete(ctx context.Context, req *pb.DeleteRequest) (*emptypb.Empty, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Создаем соединение с базой данных
	con, err := pgx.Connect(ctx, dbDSN)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer con.Close(ctx)
	// Выполняем запрос на удаление пользователя по id
	res, err := con.Exec(ctx, "DELETE FROM chats WHERE id = $1", req.Id)
	if err != nil {
		return nil, fmt.Errorf("failed to delete chat: %v", err)
	}

	// Проверяем, что удалено хотя бы одно совпадение
	if res.RowsAffected() == 0 {
		return nil, fmt.Errorf("user with id %d not found", req.Id)
	}

	// Возвращаем пустой ответ при успешном удалении
	return &emptypb.Empty{}, nil
}
func (c *ChatServer) Message(ctx context.Context, req *pb.MessageRequest) (*emptypb.Empty, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Создаем соединение с базой данных
	con, err := pgx.Connect(ctx, dbDSN)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer con.Close(ctx)

	_, err = con.Exec(ctx, "INSERT INTO message (from_user, message, timestamp) VALUES ($1, $2, $3)", req.From, req.Text, req.Timestamp.AsTime())
	if err != nil {
		return nil, fmt.Errorf("не удалось отправить сообщение: %v", err)
	}
	return &emptypb.Empty{}, nil
}
func main() {
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("Серверу пизда: %v", err)
	}
	s := grpc.NewServer()
	reflection.Register(s)
	pb.RegisterChatServer(s, &ChatServer{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Ошибка при запуске сервера: %v", err)
	}
}
