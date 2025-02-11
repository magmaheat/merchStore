package service

import "github.com/magmaheat/merchStore/internal/repo"

type StoreService struct {
	repo repo.Repo
}

func NewStoreService(repo repo.Repo) *StoreService {
	return &StoreService{repo: repo}
}
