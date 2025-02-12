package service

import (
	"github.com/magmaheat/merchStore/internal/repo"
	"time"
)

type Dependencies struct {
	Repo     repo.Repo
	SignKey  string
	TokenTTL time.Duration
}

type Service struct {
	Auth  Auth
	Store Store
}

func NewService(deps Dependencies) *Service {
	return &Service{
		Auth:  NewAuthService(deps.Repo, deps.SignKey, deps.TokenTTL),
		Store: deps.Repo,
	}
}
