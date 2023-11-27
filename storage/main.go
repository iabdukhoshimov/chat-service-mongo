package storage

import (
	"errors"
	"golang-websocker/storage/repo"
)

var ErrorTheSameId = errors.New("cannot use the same uuid for 'id' and 'parent_id' fields")
var ErrorProjectId = errors.New("not valid 'project_id'")

var (
	ChatCollection        = "chats"
	ChatMessageCollection = "chat_messages"
)

type StorageI interface {
	Chat() repo.ChatRepoI
	CloseMongoConnection()
}
