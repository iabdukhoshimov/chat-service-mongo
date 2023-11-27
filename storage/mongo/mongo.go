package mongo

import (
	"context"
	"fmt"
	"golang-websocker/config"
	"golang-websocker/storage"
	"golang-websocker/storage/repo"
	"time"

	"github.com/saidamir98/udevs_pkg/logger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Store struct {
	db       *mongo.Database
	chatRepo repo.ChatRepoI
}

func MongoConnection(log logger.LoggerI, cfg config.Config) (storage.StorageI, error) {
	log.Info("MongoDB connection...")

	url := fmt.Sprintf("mongodb://%s:%s@%s:%d/%s", cfg.MongoUser, cfg.MongoPassword, cfg.MongoHost, cfg.MongoPort, cfg.MongoDatabase)

	log.Info("url::::" + url)
	clientOptions := options.Client().
		ApplyURI(url)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, err
	}

	log.Info("MongoDB connected")

	return &Store{
		db: client.Database(cfg.MongoDatabase),
	}, nil
}

func (s *Store) CloseMongoConnection() {
	if err := s.db.Client().Disconnect(
		context.Background(),
	); err != nil {
		panic(err)
	}
}

func (s *Store) Chat() repo.ChatRepoI {
	if s.chatRepo == nil {
		s.chatRepo = NewChatRepo(s.db)
	}
	return s.chatRepo
}
