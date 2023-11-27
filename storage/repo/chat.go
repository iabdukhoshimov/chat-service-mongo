package repo

import (
	"context"
	pb "golang-websocker/genproto/chat_service"
	"golang-websocker/pkg/models"
)

// Interfaces
type ChatRepoI interface {
	CreateChat(context.Context, *pb.CreateChatRequest) (*pb.CreateChatResponse, error)
	GetChatByChatId(context.Context, *pb.GetChatByChatIdRequest) (*pb.GetChatByChatIdResponse, error)
	GetChatList(context.Context, *pb.GetChatListRequest) (*pb.GetChatListResponse, error)
	CreateMessage(context.Context, *pb.CreateMessageRequest) (*models.ChatMessage, error)
}
