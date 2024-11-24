package main

import (
	"context"
	"github.com/alexeyvas94/chat-server/api"
	pb "github.com/alexeyvas94/chat-server/proto"
	"github.com/alexeyvas94/chat-server/repository"
	"github.com/alexeyvas94/chat-server/service"
	"github.com/jackc/pgx/v4"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

const (
	dbDSN = "host=localhost port=54322 dbname=chat user=chat-user password=chat-password sslmode=disable"
)

func main() {
	// Подключаемся к базе данных
	db, err := pgx.Connect(context.Background(), dbDSN)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close(context.Background())

	// Создаем репозиторий
	chatRepo := repository.NewPostgresChatRepository(db)

	// Создаем сервис
	chatService := service.NewChatService(chatRepo)

	// Api
	chatServer := api.NewChatServer(chatService)

	// Настраиваем gRPC-сервер
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
	s := grpc.NewServer()
	reflection.Register(s)

	// Регистрируем gRPC-обработчик
	pb.RegisterChatServer(s, chatServer)

	log.Println("Server is running on port :8080")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
