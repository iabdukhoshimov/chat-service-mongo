package chat

import (
	"golang-websocker/config"
	pb "golang-websocker/genproto/chat_service"
	"golang-websocker/grpc/client"

	"golang-websocker/storage"

	"github.com/saidamir98/udevs_pkg/logger"
)

type Service struct {
	log  logger.LoggerI
	cfg  config.Config
	strg storage.StorageI
	svcs client.ServiceManagerI
	pb.UnimplementedChatServiceServer
}

func NewService(log logger.LoggerI, cfg config.Config, strg storage.StorageI, svcs client.ServiceManagerI) *Service {
	return &Service{
		log:  log,
		cfg:  cfg,
		strg: strg,
		svcs: svcs,
	}
}
