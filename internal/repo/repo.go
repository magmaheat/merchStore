package repo

import "context"

type Repo interface {
	GetUser(ctx context.Context, username string) (int, string, error)
	//CreateUser(ctx context.Context, username string, password string) (int, error)
	//CreateBalance(ctx context.Context, id int) error
	CreateUserWithBalance(ctx context.Context, username, password string) (int, error)
}
