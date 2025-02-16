package service

import (
	"github.com/magmaheat/merchStore/internal/repo"
	"github.com/magmaheat/merchStore/pkg/hasher"
	"time"
)

type Dependencies struct {
	Repo     repo.Repo
	Hasher   hasher.PasswordHasher
	SignKey  string
	TokenTTL time.Duration
}

type Service struct {
	Auth  Auth
	Store Store
}

func NewService(deps Dependencies) *Service {
	return &Service{
		Auth:  NewAuthService(deps.Repo, deps.Hasher, deps.SignKey, deps.TokenTTL),
		Store: NewStoreService(deps.Repo),
	}
}
