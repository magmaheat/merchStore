package repo

import "context"

type GetReceivedTransactionOutput struct {
	FromUser string
	Amount   int
}

type GetSentTransactionOutput struct {
	ToUser string
	Amount int
}

type Item struct {
	Type     string
	Quantity int
}

type Repo interface {
	GetUserIdWithPassword(ctx context.Context, username string) (int, string, error)
	CreateUserWithBalance(ctx context.Context, username, password string) (int, error)
	GetNameUser(ctx context.Context, userId int) (string, error)
	GetPriceItem(ctx context.Context, item string) (int, error)
	UpdateBalance(ctx context.Context, userId int, coins int) error
	AddItem(ctx context.Context, userId int, item string) error
	TransferCoins(ctx context.Context, fromUserId, toUserId, amount int) error
	GetBalance(ctx context.Context, userId int) (int, error)
	GetReceivedTransactions(ctx context.Context, userId int) ([]GetReceivedTransactionOutput, error)
	GetSentTransactions(ctx context.Context, userId int) ([]GetSentTransactionOutput, error)
	GetItems(ctx context.Context, userId int) ([]Item, error)
	//FastGetUserIdWithPassword(ctx context.Context, username, password string) (string, error)
}
