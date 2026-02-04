package server

import (
	"encoding/json"
	"fmt"
	nethttp "net/http"
	"os"
	"strings"

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
		"/auth/login":        {},
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
		if !ssoEnabled() {
			nethttp.Error(w, "sso disabled", nethttp.StatusNotFound)
			return
		}
		redirectURL := os.Getenv("SSO_SAML2_LOGIN_URL")
		if redirectURL == "" {
			redirectURL = "/auth/sso/callback?code=local-dev"
		}
		nethttp.Redirect(w, r, redirectURL, nethttp.StatusFound)
	})

	server.HandleFunc("/auth/sso/callback", func(w nethttp.ResponseWriter, r *nethttp.Request) {
		if !ssoEnabled() {
			nethttp.Error(w, "sso disabled", nethttp.StatusNotFound)
			return
		}
		code := r.URL.Query().Get("code")
		nethttp.SetCookie(w, &nethttp.Cookie{
			Name:  "sso_session",
			Value: fmt.Sprintf("session-%s", code),
			Path:  "/",
		})
		w.WriteHeader(nethttp.StatusOK)
		_, _ = w.Write([]byte("sso login success"))
	})

	server.HandleFunc("/auth/login", func(w nethttp.ResponseWriter, r *nethttp.Request) {
		if r.Method != nethttp.MethodPost {
			nethttp.Error(w, "method not allowed", nethttp.StatusMethodNotAllowed)
			return
		}

		var payload struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}

		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			nethttp.Error(w, "invalid payload", nethttp.StatusBadRequest)
			return
		}

		if strings.TrimSpace(payload.Username) == "" || strings.TrimSpace(payload.Password) == "" {
			nethttp.Error(w, "username and password required", nethttp.StatusBadRequest)
			return
		}

		nethttp.SetCookie(w, &nethttp.Cookie{
			Name:  "local_session",
			Value: fmt.Sprintf("local-%s", payload.Username),
			Path:  "/",
		})
		w.WriteHeader(nethttp.StatusOK)
		_, _ = w.Write([]byte("login success"))
	})
}

func ssoEnabled() bool {
	return strings.EqualFold(os.Getenv("SSO_ENABLED"), "true")
}
