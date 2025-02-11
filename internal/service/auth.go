package service

import (
	"github.com/magmaheat/merchStore/internal/repo"
	"time"
)

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
