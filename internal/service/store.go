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
	repo  repo.Repo
	items map[string]int
}

func NewStoreService(repo repo.Repo) *StoreService {

	return &StoreService{repo: repo,
		items: map[string]int{
			"t-shirt":    80,
			"cup":        20,
			"book":       50,
			"pen":        10,
			"powerbank":  200,
			"hoody":      300,
			"umbrella":   200,
			"socks":      10,
			"wallet":     50,
			"pink-hoody": 500,
		},
	}
}

func (s *StoreService) BuyItem(ctx context.Context, nameItem string) error {
	price, ok := s.items[nameItem]
	if !ok {
		log.Errorf("service.BuyItem: item %s not found", nameItem)
		return ErrItemNotFound
	}

	// попытаться списать с баланса стоимость

	// при удачном списаннии добавляем предмет в инвертарь

	return nil
}
