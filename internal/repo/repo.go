package repo

import "context"

type Repo interface {
	GetUserIdWithPassword(ctx context.Context, username string) (int, string, error)
	CreateUserWithBalance(ctx context.Context, username, password string) (int, error)
	GetPriceItem(ctx context.Context, item string) (int, error)
	UpdateBalance(ctx context.Context, userId int, balance int) error
	AddItem(ctx context.Context, userId int, item string) error
}
