package service

import (
	"context"
	"github.com/magmaheat/merchStore/internal/repo"
	log "github.com/sirupsen/logrus"
)

type Store interface {
	BuyItem(ctx context.Context, nameItem string) error
}

type StoreService struct {
	repo repo.Repo
}

func NewStoreService(repo repo.Repo) *StoreService {
	return &StoreService{repo: repo}
}

func (s *StoreService) BuyItem(ctx context.Context, nameItem string) error {
	price, err := s.repo.GetPriceItem(ctx, nameItem)
	if err != nil {
		return err
	}

	if price == 0 {
		return ErrItemNotFound
	}

	userId, ok := ctx.Value("userId").(int)
	log.Info("userId:", userId)
	if !ok {
		log.Errorf("service.BuyItem: %s", ErrUserIdNotFound.Error())
		return ErrUserIdNotFound
	}

	if err = s.repo.UpdateBalance(ctx, userId, -price); err != nil {
		return err
	}

	//TODO дублирует записи в таблице
	if err = s.repo.AddItem(ctx, userId, nameItem); err != nil {
		return err
	}

	return nil
}
