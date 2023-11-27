package mongo_models

type ChatMessage struct {
	ChatID     string `json:"chat_id" bson:"chat_id"`
	SenderName string `json:"sender_name" bson:"sender_name"`
	Message    string `json:"message" bson:"message"`
	Type       string `json:"type" bson:"type"`
	CreatedAt  string `json:"created_at" bson:"created_at"`
	UserID     string `json:"user_id" bson:"user_id"`
}

type Chat struct {
	Guid          string `json:"guid" bson:"guid"`
	EnvironmentID string `json:"environment_id" bson:"environment_id"`
	CreatedAt     string `json:"created_at" bson:"created_at"`
	SenderName    string `json:"sender_name" bson:"sender_name"`
	UserID        string `json:"user_id" bson:"user_id"`
}
