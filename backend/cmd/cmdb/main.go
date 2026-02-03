package main

import (
	"log"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/transport"

	"cmdb/internal/server"
)

func newApp(httpServer *server.HTTPServer, grpcServer *server.GRPCServer) *kratos.App {
	return kratos.New(
		kratos.Name("cmdb"),
		kratos.Version("0.1.0"),
		kratos.Server([]transport.Server{httpServer.Server, grpcServer.Server}...),
	)
}

func main() {
	httpServer := server.NewHTTPServer()
	grpcServer := server.NewGRPCServer()

	app := newApp(httpServer, grpcServer)
	if err := app.Run(); err != nil {
		log.Fatalf("app run failed: %v", err)
	}
}
