package grpc

import (
	"golang-websocker/config"
	"golang-websocker/genproto/chat_service"
	"golang-websocker/grpc/client"
	"golang-websocker/grpc/service/chat"
	"golang-websocker/storage"

	"github.com/saidamir98/udevs_pkg/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func SetUpServer(cfg config.Config, log logger.LoggerI, strg storage.StorageI, svcs client.ServiceManagerI) (grpcServer *grpc.Server) {

	grpcServer = grpc.NewServer()

	chat_service.RegisterChatServiceServer(grpcServer, chat.NewService(log, cfg, strg, svcs))

	reflection.Register(grpcServer)
	return
}
