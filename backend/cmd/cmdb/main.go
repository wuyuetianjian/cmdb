package main

import (
	"context"
	"log"
	"os"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/transport"

	"cmdb/internal/data"
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
	ctx := context.Background()
	driver := os.Getenv("DB_DRIVER")
	if driver == "" {
		driver = "sqlite"
	}
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		dsn = "file:cmdb.db?_fk=1"
	}

	entClient, err := data.NewEntClient(ctx, driver, dsn)
	if err != nil {
		log.Fatalf("init ent client failed: %v", err)
	}

	httpServer := server.NewHTTPServer(entClient)
	grpcServer := server.NewGRPCServer()

	app := newApp(httpServer, grpcServer)
	if err := app.Run(); err != nil {
		log.Fatalf("app run failed: %v", err)
	}
}
