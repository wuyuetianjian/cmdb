package server

import (
	"github.com/go-kratos/kratos/v2/transport/grpc"
)

type GRPCServer struct {
	*grpc.Server
}

func NewGRPCServer() *GRPCServer {
	server := grpc.NewServer(
		grpc.Address(":9000"),
	)

	return &GRPCServer{Server: server}
}
