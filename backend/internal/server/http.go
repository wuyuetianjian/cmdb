package server

import (
	"github.com/go-kratos/kratos/v2/transport/http"
)

type HTTPServer struct {
	*http.Server
}

func NewHTTPServer() *HTTPServer {
	server := http.NewServer(
		http.Address(":8000"),
	)

	return &HTTPServer{Server: server}
}
