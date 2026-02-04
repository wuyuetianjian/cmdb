package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"os"
	"sync"

	"golang.org/x/crypto/bcrypt"
)

const (
	DefaultAdminUsername = "admin"
	DefaultAdminPassword = "admin123"
)

type User struct {
	Username           string
	PasswordHash       string
	Tags               []string
	IsAdmin            bool
	MustChangePassword bool
}

type Session struct {
	Username string
	Tags     []string
	IsAdmin  bool
}

type Manager struct {
	mu       sync.RWMutex
	users    map[string]*User
	sessions map[string]Session
}

func NewManager() *Manager {
	manager := &Manager{
		users:    make(map[string]*User),
		sessions: make(map[string]Session),
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
}

func (m *Manager) Authenticate(username, password string) (string, *User, error) {
	m.mu.Lock()
	m.ensureDefaultAdmin()
	m.mu.Unlock()

	m.mu.RLock()
	user, ok := m.users[username]
	m.mu.RUnlock()
	if !ok {
		return "", nil, errors.New("invalid credentials")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
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
	return nil
}

func (m *Manager) Session(sessionID string) (Session, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	session, ok := m.sessions[sessionID]
	return session, ok
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

func newSessionID() string {
	var buf [16]byte
	_, _ = rand.Read(buf[:])
	return hex.EncodeToString(buf[:])
}
