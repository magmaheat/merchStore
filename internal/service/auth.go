package service

import (
	"context"
	"github.com/magmaheat/merchStore/internal/repo"
	"time"
)

type AuthGenerateTokenInput struct {
	Username string
	Password string
}

type Auth interface {
	GenerateToken(ctx context.Context, input AuthGenerateTokenInput) (string, error)
}

type AuthService struct {
	repo     repo.Repo
	signKey  string
	tokenTTL time.Duration
}

func NewAuthService(repo repo.Repo, signKey string, tokenTTL time.Duration) *AuthService {
	return &AuthService{
		repo:     repo,
		signKey:  signKey,
		tokenTTL: tokenTTL,
	}
}

func (s *AuthService) GenerateToken(ctx context.Context, input AuthGenerateTokenInput) (string, error) {
	return "", nil
}
