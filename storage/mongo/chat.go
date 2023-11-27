package mongo

import (
	"context"
	pb "golang-websocker/genproto/chat_service"
	"golang-websocker/pkg/models"
	"golang-websocker/storage"
	"golang-websocker/storage/repo"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"
)

type ChatRepo struct {
	db *mongo.Database
}

func NewChatRepo(db *mongo.Database) repo.ChatRepoI { //!!!
	return &ChatRepo{
		db: db,
	}
}

func (self *ChatRepo) CreateChat(ctx context.Context, req *pb.CreateChatRequest) (resp *pb.CreateChatResponse, err error) {

	resp = &pb.CreateChatResponse{}
	id := uuid.New().String()

	chat := &models.CreateChat{
		UserID:        req.UserId,
		Guid:          id,
		SenderName:    req.Chat.SenderName,
		EnvironmentID: req.Chat.EnvironmentId,
		CreatedAt:     time.Now().Format(time.RFC3339),
	}
	Chat := self.db.Collection(storage.ChatCollection)

	_, err = Chat.InsertOne(ctx, chat)
	if err != nil {
		return resp, err
	}
	resp.ChatId = id
	return resp, nil
}

func (self *ChatRepo) GetChatByChatId(ctx context.Context, req *pb.GetChatByChatIdRequest) (resp *pb.GetChatByChatIdResponse, err error) {
	resp = &pb.GetChatByChatIdResponse{}
	chat := &models.Chat{}

	var (
		userMessage []*pb.UserMessage
	)

	pipline := mongo.Pipeline{
		bson.D{{Key: "$match", Value: bson.D{
			{Key: "guid", Value: req.GetChatId()},
		}}},
		bson.D{
			{Key: "$project", Value: bson.D{
				{Key: "guid", Value: 1},
				{Key: "environment_id", Value: 1},
				{Key: "created_at", Value: 1},
			}},
		},
		bson.D{
			{Key: "$lookup", Value: bson.D{
				{Key: "from", Value: storage.ChatMessageCollection},
				{Key: "localField", Value: "guid"},
				{Key: "foreignField", Value: "chat_id"},
				{Key: "as", Value: "messages"},
			}},
		},
		bson.D{
			{Key: "$sort", Value: bson.D{
				{Key: "created_at", Value: -1},
			}},
		},
		bson.D{
			{Key: "$limit", Value: 30},
		},
	}

	cur, err := self.db.Collection(storage.ChatCollection).Aggregate(ctx, pipline)
	if err != nil {
		return resp, err
	}

	for cur.Next(ctx) {
		err := cur.Decode(chat)
		if err != nil {
			return resp, err
		}
	}

	for _, message := range chat.Messages {
		userMessage = append(userMessage, &pb.UserMessage{
			SenderName: message.SenderName,
			Message:    message.Message,
			CreatedAt:  message.CreatedAt,
			Type:       message.Type,
			UserId:     message.UserID,
		})
	}

	return &pb.GetChatByChatIdResponse{
		ChatId:   chat.Guid,
		Messages: userMessage,
	}, nil
}

func (self *ChatRepo) GetChatList(ctx context.Context, req *pb.GetChatListRequest) (resp *pb.GetChatListResponse, err error) {
	resp = &pb.GetChatListResponse{}

	chat := &models.ChatContent{}

	var (
		chats []*pb.ChatWithLastMessageData
	)

	pipline := mongo.Pipeline{
		bson.D{
			{Key: "$match", Value: bson.D{
				{Key: "environment_id", Value: req.GetEnvironmentId()},
			}},
		},
		bson.D{
			{Key: "$lookup", Value: bson.D{
				{Key: "from", Value: storage.ChatMessageCollection},
				{Key: "localField", Value: "guid"},
				{Key: "foreignField", Value: "chat_id"},
				{Key: "as", Value: "message"},
			}},
		},
		bson.D{
			{Key: "$unwind", Value: bson.D{
				{Key: "path", Value: "$message"},
				{Key: "preserveNullAndEmptyArrays", Value: true},
			}},
		},
		bson.D{
			{Key: "$sort", Value: bson.D{
				{Key: "message.created_at", Value: -1},
			}},
		},
		bson.D{
			{Key: "$group", Value: bson.D{
				{Key: "_id", Value: "$guid"},
				{Key: "guid", Value: bson.D{
					{Key: "$first", Value: "$guid"},
				}},
				{Key: "environment_id", Value: bson.D{
					{Key: "$first", Value: "$environment_id"},
				}},
				{Key: "created_at", Value: bson.D{
					{Key: "$first", Value: "$created_at"},
				}},
				{Key: "sender_name", Value: bson.D{
					{Key: "$first", Value: "$sender_name"},
				}},
				{Key: "message", Value: bson.D{
					{Key: "$first", Value: "$message"},
				}},
			}},
		},
		bson.D{
			{Key: "$project", Value: bson.D{
				{Key: "guid", Value: 1},
				{Key: "environment_id", Value: 1},
				{Key: "created_at", Value: 1},
				{Key: "sender_name", Value: 1},
				{Key: "message", Value: bson.D{
					{Key: "$ifNull", Value: bson.A{"$message", bson.D{}}},
				}},
			}},
		},
	}

	// If search parameter is provided, add it to the pipeline
	if req.GetSearch() != "" {
		regex := primitive.Regex{Pattern: req.GetSearch(), Options: "i"}
		pipline = append(pipline, bson.D{
			{Key: "$match", Value: bson.D{
				{Key: "$or", Value: bson.A{
					bson.D{
						{Key: "sender_name", Value: regex},
					},
				}},
			}},
		})
	}

	cur, err := self.db.Collection(storage.ChatCollection).Aggregate(ctx, pipline)
	if err != nil {
		return resp, err
	}
	for cur.Next(ctx) {
		err := cur.Decode(chat)
		if err != nil {
			return resp, err
		}

		chats = append(chats, &pb.ChatWithLastMessageData{
			ChatId: chat.Guid,
			Message: &pb.UserMessage{
				SenderName: chat.Message.SenderName,
				Message:    chat.Message.Message,
				CreatedAt:  chat.Message.CreatedAt,
				Type:       chat.Message.Type,
				UserId:     chat.Message.UserID,
			},
			SenderName: chat.SenderName,
			CreatedAt:  chat.CreatedAt,
		})
	}
	return &pb.GetChatListResponse{
		EnvironmentId: req.GetEnvironmentId(),
		Chats:         chats,
	}, nil
}

func (self *ChatRepo) CreateMessage(ctx context.Context, req *pb.CreateMessageRequest) (messaage *models.ChatMessage, err error) {

	chatMessage := &models.ChatMessage{
		ChatID:     req.GetChatId(),
		CreatedAt:  time.Now().Format(time.RFC3339),
		UserID:     req.GetMessage().GetUserId(),
		SenderName: req.GetMessage().GetSenderName(),
		Message:    req.GetMessage().GetMessage(),
		Type:       req.GetMessage().GetType(),
	}

	_, err = self.db.Collection(storage.ChatMessageCollection).InsertOne(ctx, chatMessage)
	if err != nil {
		return nil, err
	}

	return chatMessage, nil
}
