//go:build ent

package data

import (
	"context"
	"errors"
	"fmt"

	"cmdb/internal/biz"
	"cmdb/internal/data/ent"
	"cmdb/internal/data/ent/user"

	"golang.org/x/crypto/bcrypt"
)

type AuthRepository struct {
	ent *ent.Client
}

func NewAuthRepository(entClient *ent.Client) *AuthRepository {
	return &AuthRepository{ent: entClient}
}

func (r *AuthRepository) CreateUser(ctx context.Context, username, password string) (*biz.User, error) {
	hashed, err := hashPassword(password)
	if err != nil {
		return nil, err
	}
	entity, err := r.ent.User.
		Create().
		SetUsername(username).
		SetPasswordHash(hashed).
		Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}
	return &biz.User{
		ID:        int64(entity.ID),
		Username:  entity.Username,
		CreatedAt: entity.CreatedAt.Format("2006-01-02 15:04:05"),
	}, nil
}

func (r *AuthRepository) Authenticate(ctx context.Context, username, password string) (*biz.User, error) {
	entity, err := r.ent.User.
		Query().
		Where(user.Username(username)).
		Only(ctx)
	if err != nil {
		return nil, fmt.Errorf("query user: %w", err)
	}
	if err := verifyPassword(entity.PasswordHash, password); err != nil {
		return nil, errors.New("invalid credentials")
	}
	return &biz.User{
		ID:        int64(entity.ID),
		Username:  entity.Username,
		CreatedAt: entity.CreatedAt.Format("2006-01-02 15:04:05"),
	}, nil
}

func hashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("hash password: %w", err)
	}
	return string(hashed), nil
}

func verifyPassword(hashedPassword string, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
