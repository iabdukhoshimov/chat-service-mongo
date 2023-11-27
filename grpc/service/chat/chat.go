package chat

import (
	"context"
	"fmt"
	pb "golang-websocker/genproto/chat_service"

	"github.com/saidamir98/udevs_pkg/logger"
)

func (self *Service) CreateChat(ctx context.Context, req *pb.CreateChatRequest) (resp *pb.CreateChatResponse, err error) {
	self.log.Info("---CreateChat---", logger.Any("req", req))
	fmt.Println(req)
	chat, err := self.strg.Chat().CreateChat(ctx, req)
	if err != nil {
		self.log.Error("---Err->CreateChat-->", logger.Error(err))
		return nil, err
	}

	resp = &pb.CreateChatResponse{
		ChatId: chat.ChatId,
	}
	return resp, nil
}

func (self *Service) GetChatList(ctx context.Context, req *pb.GetChatListRequest) (resp *pb.GetChatListResponse, err error) {
	self.log.Info("---GetChatList---", logger.Any("req", req))
	chatlist, err := self.strg.Chat().GetChatList(ctx, req)
	if err != nil {
		self.log.Error("---Err->GetChatList-->", logger.Error(err))
		return nil, err
	}

	return chatlist, nil
}

func (self *Service) GetChatByChatId(ctx context.Context, req *pb.GetChatByChatIdRequest) (*pb.GetChatByChatIdResponse, error) {
	self.log.Info("---GetChatByChatId---", logger.Any("req", req))

	chat, err := self.strg.Chat().GetChatByChatId(ctx, req)
	if err != nil {
		self.log.Error("---Err->GetChatByChatId-->", logger.Error(err))
		return nil, err
	}

	return chat, nil
}

// func (self *Service) CreateMessage(ctx context.Context, req *pb.CreateMessageRequest) (resp *emptypb.Empty, err error) {
// 	self.log.Info("---CreateMessage---", logger.Any("req", req))

// 	resp = &emptypb.Empty{}
// 	err = self.strg.Chat().CreateMessage(ctx, req)
// 	if err != nil {
// 		self.log.Error("---Err->CreateMessage-->", logger.Error(err))
// 		return resp, err
// 	}

// 	return resp, nil
// }
