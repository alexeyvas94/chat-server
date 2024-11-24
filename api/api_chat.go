package api

import (
	"context"
	pb "github.com/alexeyvas94/chat-server/proto"
	"github.com/alexeyvas94/chat-server/service"
	"google.golang.org/protobuf/types/known/emptypb"
	"sync"
)

type ChatServer struct {
	pb.UnimplementedChatServer
	mu          sync.Mutex
	chatService service.ChatService
}

func NewChatServer(chatService service.ChatService) *ChatServer {
	return &ChatServer{chatService: chatService}
}

func (c *ChatServer) Create(ctx context.Context, req *pb.CreateRequest) (*pb.CreateResponse, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	chatID, err := c.chatService.CreateChat(ctx, req.Usernames)
	if err != nil {
		return nil, err
	}
	return &pb.CreateResponse{Id: chatID}, nil
}

func (c *ChatServer) Delete(ctx context.Context, req *pb.DeleteRequest) (*emptypb.Empty, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if err := c.chatService.DeleteChat(ctx, req.Id); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (c *ChatServer) Message(ctx context.Context, req *pb.MessageRequest) (*emptypb.Empty, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if err := c.chatService.AddMessage(ctx, req.From, req.Text, req.Timestamp.AsTime()); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
