package models

type Chat struct {
	Guid          string        `json:"guid" bson:"guid"`
	Messages      []ChatMessage `json:"messages" bson:"messages"`
	EnvironmentID string        `json:"environment_id" bson:"environment_id"`
	CreatedAt     string        `json:"created_at" bson:"created_at"`
}

type ChatMessage struct {
	SenderName string `json:"sender_name" bson:"sender_name"`
	Message    string `json:"message" bson:"message"`
	CreatedAt  string `json:"created_at" bson:"created_at"`
	Type       string `json:"type" bson:"type"`
	ChatID     string `json:"chat_id" bson:"chat_id"`
	UserID     string `json:"user_id" bson:"user_id"`
}

type CreateChat struct {
	Guid          string `json:"guid" bson:"guid"`
	SenderName    string `json:"sender_name" bson:"sender_name"`
	CreatedAt     string `json:"created_at" bson:"created_at"`
	EnvironmentID string `json:"environment_id" bson:"environment_id"`
	UserID        string `json:"user_id" bson:"user_id"`
}

type ChatContent struct {
	Guid          string      `json:"guid" bson:"guid"`
	Message       ChatMessage `json:"message" bson:"message"`
	EnvironmentID string      `json:"environment_id" bson:"environment_id"`
	CreatedAt     string      `json:"created_at" bson:"created_at"`
	SenderName    string      `json:"sender_name" bson:"sender_name"`
}

type Message struct {
	Message string `json:"message" bson:"message"`
	Type    string `json:"type" bson:"type"`
}

type CtxData struct {
	UserID     string `json:"user_id"`
	SenderName string `json:"sender_name"`
}