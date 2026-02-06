package auth

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"

	"golang.org/x/crypto/bcrypt"
)

const (
	DefaultAdminUsername = "admin"
	DefaultAdminPassword = "admin123"
)

type User struct {
	Username           string   `json:"username"`
	PasswordHash       string   `json:"password_hash"`
	Tags               []string `json:"tags"`
	IsAdmin            bool     `json:"is_admin"`
	MustChangePassword bool     `json:"must_change_password"`
}

type Session struct {
	Username string
	Tags     []string
	IsAdmin  bool
}

type Manager struct {
	mu        sync.RWMutex
	users     map[string]*User
	sessions  map[string]Session
	usersFile string
}

func NewManager() *Manager {
	manager := &Manager{
		users:     make(map[string]*User),
		sessions:  make(map[string]Session),
		usersFile: defaultUsersFile(),
	}
	if err := manager.loadUsers(); err != nil {
		log.Printf("auth: load users failed, fallback to empty store: %v", err)
	}
	manager.ensureDefaultAdmin()
	return manager
}

func (m *Manager) ensureDefaultAdmin() {
	if _, ok := m.users[DefaultAdminUsername]; ok {
		return
	}
	hash, _ := bcrypt.GenerateFromPassword([]byte(defaultAdminPassword()), bcrypt.DefaultCost)
	m.users[DefaultAdminUsername] = &User{
		Username:           DefaultAdminUsername,
		PasswordHash:       string(hash),
		Tags:               []string{"assets:read", "assets:write", "admin:write"},
		IsAdmin:            true,
		MustChangePassword: true,
	}
	if err := m.saveUsers(); err != nil {
		log.Printf("auth: persist default admin failed: %v", err)
	}
}

func (m *Manager) Authenticate(username, password string) (string, *User, error) {
	m.mu.Lock()
	m.ensureDefaultAdmin()
	m.mu.Unlock()

	m.mu.RLock()
	user, ok := m.users[username]
	m.mu.RUnlock()
	if !ok {
		log.Printf("auth: user not found (%s)", username)
		return "", nil, errors.New("invalid credentials")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		log.Printf("auth: password mismatch (%s): %v", username, err)
		return "", nil, errors.New("invalid credentials")
	}
	sessionID := newSessionID()
	m.mu.Lock()
	m.sessions[sessionID] = Session{Username: user.Username, Tags: user.Tags, IsAdmin: user.IsAdmin}
	m.mu.Unlock()
	return sessionID, user, nil
}

func defaultAdminPassword() string {
	if value := os.Getenv("DEFAULT_ADMIN_PASSWORD"); value != "" {
		return value
	}
	return DefaultAdminPassword
}

func defaultUsersFile() string {
	if value := os.Getenv("USER_DB_FILE"); value != "" {
		return value
	}
	return "users.db.json"
}

func (m *Manager) RegisterSSOSession(username string, tags []string, isAdmin bool) string {
	sessionID := newSessionID()
	m.mu.Lock()
	m.sessions[sessionID] = Session{Username: username, Tags: tags, IsAdmin: isAdmin}
	m.mu.Unlock()
	return sessionID
}

func (m *Manager) ChangePassword(username, newPassword string) error {
	if newPassword == "" {
		return errors.New("password required")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	user, ok := m.users[username]
	if !ok {
		return errors.New("user not found")
	}
	user.PasswordHash = string(hash)
	user.MustChangePassword = false
	if err := m.saveUsers(); err != nil {
		return err
	}
	return nil
}

func (m *Manager) Session(sessionID string) (Session, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	session, ok := m.sessions[sessionID]
	return session, ok
}

func (m *Manager) IsAdmin(sessionID string) bool {
	session, ok := m.Session(sessionID)
	if !ok {
		return false
	}
	return session.IsAdmin
}

func (m *Manager) RequirePermission(sessionID, permission string) bool {
	session, ok := m.Session(sessionID)
	if !ok {
		return false
	}
	if session.IsAdmin {
		return true
	}
	for _, tag := range session.Tags {
		if tag == permission {
			return true
		}
	}
	return false
}

func (m *Manager) loadUsers() error {
	content, err := os.ReadFile(m.usersFile)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return fmt.Errorf("read users file: %w", err)
	}
	if len(content) == 0 {
		return nil
	}
	users := make(map[string]*User)
	if err := json.Unmarshal(content, &users); err != nil {
		return fmt.Errorf("decode users file: %w", err)
	}
	m.users = users
	return nil
}

func (m *Manager) saveUsers() error {
	if err := os.MkdirAll(filepath.Dir(m.usersFile), 0o755); err != nil && filepath.Dir(m.usersFile) != "." {
		return fmt.Errorf("mkdir users dir: %w", err)
	}
	content, err := json.MarshalIndent(m.users, "", "  ")
	if err != nil {
		return fmt.Errorf("encode users file: %w", err)
	}
	if err := os.WriteFile(m.usersFile, content, 0o600); err != nil {
		return fmt.Errorf("write users file: %w", err)
	}
	return nil
}

func newSessionID() string {
	var buf [16]byte
	_, _ = rand.Read(buf[:])
	return hex.EncodeToString(buf[:])
}
