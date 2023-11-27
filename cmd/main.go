package main

import (
	"net"

	"golang-websocker/api/handlers"
	"golang-websocker/config"
	"golang-websocker/grpc"
	"golang-websocker/grpc/client"
	"golang-websocker/pkg/wsclient"
	"golang-websocker/storage/mongo"

	"golang-websocker/api"

	"github.com/gin-gonic/gin"
	"github.com/saidamir98/udevs_pkg/logger"
)

func main() {
	cfg := config.Load()

	var loggerLevel string

	switch cfg.Environment {
	case config.DebugMode:
		loggerLevel = logger.LevelDebug
		gin.SetMode(gin.DebugMode)
	case config.TestMode:
		loggerLevel = logger.LevelDebug
		gin.SetMode(gin.TestMode)
	default:
		loggerLevel = logger.LevelInfo
		gin.SetMode(gin.ReleaseMode)
	}
	log := logger.NewLogger(cfg.ServiceName, loggerLevel)
	defer func() { _ = logger.Cleanup(log) }()

	// Connect to postgres
	store, err := mongo.MongoConnection(log, cfg)
	if err != nil {
		log.Panic("mongo.MongoConnection", logger.Error(err))
	}
	defer store.CloseMongoConnection()

	// Connect to grpc clients
	svcs, err := client.NewGrpcClients(cfg)
	if err != nil {
		log.Panic("client.NewGrpcClients", logger.Error(err))
	}

	// Run grpc server
	grpcServer := grpc.SetUpServer(cfg, log, store, svcs)

	go func() {
		lis, err := net.Listen("tcp", cfg.RPCPort)
		if err != nil {
			log.Panic("net.Listen", logger.Error(err))
		}

		log.Info("GRPC: Server being started...", logger.String("port", cfg.RPCPort))

		if err := grpcServer.Serve(lis); err != nil {
			log.Panic("grpcServer.Serve", logger.Error(err))
		}
	}()

	hub := wsclient.NewHub(store)
	go hub.Run()

	h := handlers.NewHandler(cfg, log, svcs, store, hub)

	r := api.SetUpRouter(h, cfg)

	if err := r.Run(cfg.HTTPPort); err != nil {
		log.Panic("r.Run", logger.Error(err))
	}
}
