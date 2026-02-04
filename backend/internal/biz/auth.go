package biz

import (
	"context"
)

type User struct {
	ID        int64
	Username  string
	CreatedAt string
}

type AuthRepo interface {
	CreateUser(ctx context.Context, username, password string) (*User, error)
	Authenticate(ctx context.Context, username, password string) (*User, error)
}

type AuthUsecase struct {
	repo AuthRepo
}

func NewAuthUsecase(repo AuthRepo) *AuthUsecase {
	return &AuthUsecase{repo: repo}
}

func (uc *AuthUsecase) Register(ctx context.Context, username, password string) (*User, error) {
	return uc.repo.CreateUser(ctx, username, password)
}

func (uc *AuthUsecase) Login(ctx context.Context, username, password string) (*User, error) {
	return uc.repo.Authenticate(ctx, username, password)
}
