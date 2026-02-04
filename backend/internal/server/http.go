package server

import (
	"encoding/json"
	"fmt"
	nethttp "net/http"
	"os"
	"strings"

	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/go-kratos/kratos/v2/transport/http"

	"cmdb/internal/auth"
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
			internalMiddleware.AuthMiddleware(publicPaths, validateRequest),
		),
	)

	registerPublicRoutes(server)

	return &HTTPServer{Server: server}
}

var authManager = auth.NewManager()

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
		sessionID := authManager.RegisterSSOSession("sso-user", []string{"assets:read"}, false)
		nethttp.SetCookie(w, &nethttp.Cookie{
			Name:  "sso_session",
			Value: fmt.Sprintf("%s-%s", sessionID, code),
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

		sessionID, user, err := authManager.Authenticate(payload.Username, payload.Password)
		if err != nil {
			nethttp.Error(w, "invalid credentials: check username/password", nethttp.StatusUnauthorized)
			return
		}

		nethttp.SetCookie(w, &nethttp.Cookie{
			Name:  "local_session",
			Value: sessionID,
			Path:  "/",
		})

		response := struct {
			Username           string   `json:"username"`
			MustChangePassword bool     `json:"must_change_password"`
			Tags               []string `json:"tags"`
		}{
			Username:           user.Username,
			MustChangePassword: user.MustChangePassword,
			Tags:               user.Tags,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(nethttp.StatusOK)
		_ = json.NewEncoder(w).Encode(response)
	})

	server.HandleFunc("/auth/password", func(w nethttp.ResponseWriter, r *nethttp.Request) {
		if r.Method != nethttp.MethodPost {
			nethttp.Error(w, "method not allowed", nethttp.StatusMethodNotAllowed)
			return
		}

		sessionID, err := extractSessionID(r)
		if err != nil {
			nethttp.Error(w, "missing session", nethttp.StatusUnauthorized)
			return
		}

		var payload struct {
			NewPassword string `json:"new_password"`
		}

		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			nethttp.Error(w, "invalid payload", nethttp.StatusBadRequest)
			return
		}

		session, ok := authManager.Session(sessionID)
		if !ok {
			nethttp.Error(w, "invalid session", nethttp.StatusUnauthorized)
			return
		}

		if err := authManager.ChangePassword(session.Username, payload.NewPassword); err != nil {
			nethttp.Error(w, err.Error(), nethttp.StatusBadRequest)
			return
		}

		w.WriteHeader(nethttp.StatusOK)
		_, _ = w.Write([]byte("password updated"))
	})

	server.HandleFunc("/assets", func(w nethttp.ResponseWriter, r *nethttp.Request) {
		sessionID, err := extractSessionID(r)
		if err != nil {
			nethttp.Error(w, "missing session", nethttp.StatusUnauthorized)
			return
		}

		required := "assets:read"
		if r.Method == nethttp.MethodPost {
			required = "assets:write"
		}
		if !authManager.RequirePermission(sessionID, required) {
			nethttp.Error(w, "forbidden", nethttp.StatusForbidden)
			return
		}

		w.WriteHeader(nethttp.StatusOK)
		_, _ = w.Write([]byte("assets ok"))
	})
}

func ssoEnabled() bool {
	return strings.EqualFold(os.Getenv("SSO_ENABLED"), "true")
}

func validateRequest(operation string, header transport.Header) error {
	cookie := header.Get("Cookie")
	if cookie == "" {
		return fmt.Errorf("missing session")
	}
	sessionID := extractSessionIDFromCookie(cookie)
	if sessionID == "" {
		return fmt.Errorf("missing session")
	}

	if strings.HasPrefix(operation, "/assets") {
		if !authManager.RequirePermission(sessionID, "assets:read") {
			return fmt.Errorf("permission denied")
		}
	}

	return nil
}

func extractSessionID(r *nethttp.Request) (string, error) {
	cookie, err := r.Cookie("local_session")
	if err == nil && cookie.Value != "" {
		return cookie.Value, nil
	}
	cookie, err = r.Cookie("sso_session")
	if err == nil && cookie.Value != "" {
		return strings.Split(cookie.Value, "-")[0], nil
	}
	return "", fmt.Errorf("missing session")
}

func extractSessionIDFromCookie(cookieHeader string) string {
	for _, item := range strings.Split(cookieHeader, ";") {
		item = strings.TrimSpace(item)
		if strings.HasPrefix(item, "local_session=") {
			return strings.TrimPrefix(item, "local_session=")
		}
		if strings.HasPrefix(item, "sso_session=") {
			value := strings.TrimPrefix(item, "sso_session=")
			return strings.Split(value, "-")[0]
		}
	}
	return ""
}
