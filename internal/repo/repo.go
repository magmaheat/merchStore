package repo

import "context"

type Repo interface {
	GetUserIdWithPassword(ctx context.Context, username string) (int, string, error)
	CreateUserWithBalance(ctx context.Context, username, password string) (int, error)
	GetNameUser(ctx context.Context, userId int) (string, error)
	GetPriceItem(ctx context.Context, item string) (int, error)
	UpdateBalance(ctx context.Context, userId int, coins int) error
	AddItem(ctx context.Context, userId int, item string) error
	TransferCoins(ctx context.Context, fromUserId, toUserId, amount int) error
}
