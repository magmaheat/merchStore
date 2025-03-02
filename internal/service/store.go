package service

import (
	"context"
	"github.com/magmaheat/merchStore/internal/models"
	"github.com/magmaheat/merchStore/internal/repo"
	"github.com/magmaheat/merchStore/internal/types"
	log "github.com/sirupsen/logrus"
)

//go:generate go run github.com/vektra/mockery/v2@latest --name=Store --output=./../mocks
type Store interface {
	BuyItem(ctx context.Context, nameItem string) error
	SendCoin(ctx context.Context, toUser string, amount int) error
	GetInfo(ctx context.Context) (*models.Info, error)
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

	userId, ok := ctx.Value(types.UserIdCtx).(int)
	if !ok {
		log.Errorf("service.Store.BuyItem: %s", ErrUserIdNotFound.Error())
		return ErrUserIdNotFound
	}

	if err = s.repo.UpdateBalance(ctx, userId, -price); err != nil {
		return err
	}

	if err = s.repo.AddItem(ctx, userId, nameItem); err != nil {
		return err
	}

	return nil
}

func (s *StoreService) SendCoin(ctx context.Context, toUser string, amount int) error {
	toUserId, _, err := s.repo.GetUserIdWithPassword(ctx, toUser)
	if err != nil {
		return err
	}

	fromUserId, ok := ctx.Value(types.UserIdCtx).(int)
	if !ok {
		log.Errorf("service.Store.SendCoin: %s", ErrUserIdNotFound)
		return ErrUserIdNotFound
	}

	err = s.repo.TransferCoins(ctx, fromUserId, toUserId, amount)
	if err != nil {
		return err
	}

	return nil
}

func (s *StoreService) GetInfo(ctx context.Context) (*models.Info, error) {
	userId, ok := ctx.Value(types.UserIdCtx).(int)
	if !ok {
		log.Errorf("service.Store.GetInfo: %s", ErrUserIdNotFound)
		return nil, ErrUserIdNotFound
	}

	coins, err := s.repo.GetBalance(ctx, userId)
	if err != nil {
		return nil, err
	}

	items, err := s.repo.GetItems(ctx, userId)
	if err != nil {
		return nil, err
	}

	receiverTransactions, err := s.repo.GetReceivedTransactions(ctx, userId)
	if err != nil {
		return nil, err
	}

	senderTransactions, err := s.repo.GetSentTransactions(ctx, userId)
	if err != nil {
		return nil, err
	}

	receivedCoins := models.ConvertReceivedTransactions(receiverTransactions)
	sentCoins := models.ConvertSentTransactions(senderTransactions)
	inventory := models.ConvertInventory(items)

	return models.NewInfo(
		coins,
		inventory,
		receivedCoins,
		sentCoins,
	), nil

}
