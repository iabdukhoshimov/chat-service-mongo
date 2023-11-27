package client

import (
	"golang-websocker/config"
)

// Client is a gRPC client for the service.

type ServiceManagerI interface {
}

type grpcClients struct {
}

func NewGrpcClients(cfg config.Config) (ServiceManagerI, error) {

	return &grpcClients{}, nil
}
