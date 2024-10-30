package main

import (
	"context"
	"errors"
	pb "github.com/alexeyvas94/chat-server/pkg"
	"github.com/brianvoe/gofakeit/v6"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
	"net"
	"sync"
	"time"
)

type Text struct {
	name_user []string
	chat_text string
}

type ChatServer struct {
	pb.UnimplementedChatServer
	list_chat map[int64]Text
	mu        sync.Mutex
}

func ConvertTimestampToString(ts *timestamppb.Timestamp) (string, error) {
	// Преобразуем в time.Time
	t := ts.AsTime()

	// Форматируем в строку
	return t.Format(time.RFC3339), nil
}

func (c *ChatServer) Create(ctx context.Context, req *pb.CreateRequest) (*pb.CreateResponse, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	id := gofakeit.Number(1, 1000)
	c.list_chat[int64(id)] = Text{name_user: req.Usernames}
	return &pb.CreateResponse{Id: int64(id)}, nil
}
func (c *ChatServer) Delete(ctx context.Context, req *pb.DeleteRequest) (*emptypb.Empty, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if _, ok := c.list_chat[int64(req.Id)]; ok {
		delete(c.list_chat, int64(req.Id))
	} else {
		errors.New("Чат не найден")
	}
	return &emptypb.Empty{}, nil
}
func (c *ChatServer) Message(ctx context.Context, req *pb.MessageRequest) (*emptypb.Empty, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	chat, exists := c.list_chat[int64(req.Id)]
	if !exists {
		return nil, errors.New("Чат не найден")
	} else {
		s, _ := ConvertTimestampToString(req.Timestamp)
		chat.chat_text = chat.chat_text + req.From + s + req.Text
	}
	log.Println(chat.chat_text)
	return &emptypb.Empty{}, nil
}
func main() {
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("Серверу пизда: %v", err)
	}
	s := grpc.NewServer()
	reflection.Register(s)
	pb.RegisterChatServer(s, &ChatServer{
		list_chat: make(map[int64]Text),
	})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Ошибка при запуске сервера: %v", err)
	}
}
