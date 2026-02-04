package server

import (
	"fmt"
	nethttp "net/http"
	"os"

	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/http"

	internalMiddleware "cmdb/internal/middleware"
)

type HTTPServer struct {
	*http.Server
}

func NewHTTPServer() *HTTPServer {
	publicPaths := map[string]struct{}{
		"/health":            {},
		"/auth/sso/login":    {},
		"/auth/sso/callback": {},
	}

	server := http.NewServer(
		http.Address(":8000"),
		http.Middleware(
			recovery.Recovery(),
			internalMiddleware.AuthMiddleware(publicPaths),
		),
	)

	registerPublicRoutes(server)

	return &HTTPServer{Server: server}
}

func registerPublicRoutes(server *http.Server) {
	server.HandleFunc("/health", func(w nethttp.ResponseWriter, r *nethttp.Request) {
		w.WriteHeader(nethttp.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	server.HandleFunc("/auth/sso/login", func(w nethttp.ResponseWriter, r *nethttp.Request) {
		redirectURL := os.Getenv("SSO_REDIRECT_URL")
		if redirectURL == "" {
			redirectURL = "/auth/sso/callback?code=local-dev"
		}
		nethttp.Redirect(w, r, redirectURL, nethttp.StatusFound)
	})

	server.HandleFunc("/auth/sso/callback", func(w nethttp.ResponseWriter, r *nethttp.Request) {
		code := r.URL.Query().Get("code")
		nethttp.SetCookie(w, &nethttp.Cookie{
			Name:  "sso_session",
			Value: fmt.Sprintf("session-%s", code),
			Path:  "/",
		})
		w.WriteHeader(nethttp.StatusOK)
		_, _ = w.Write([]byte("sso login success"))
	})
}
