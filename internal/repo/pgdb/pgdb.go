package pgdb

import (
	"context"
	"errors"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/magmaheat/merchStore/internal/repo"
	"github.com/magmaheat/merchStore/pkg/postgres"
	log "github.com/sirupsen/logrus"
	"time"
)

type Storage struct {
	db *postgres.Postgres
}

func New(db *postgres.Postgres) *Storage {
	return &Storage{db: db}
}

func (s *Storage) GetUserIdWithPassword(ctx context.Context, username string) (int, string, error) {
	sql, args, _ := s.db.Builder.
		Select("id", "password").
		From("users").
		Where("username = ?", username).
		ToSql()

	var id int
	var hash string

	err := s.db.Pool.QueryRow(ctx, sql, args...).Scan(&id, &hash)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, "", nil
		}
		log.Errorf("repo.pgdb.GetUserIdWithPassword.QueryRow: %v:", err)
		return 0, "", fmt.Errorf("could not get user: %w", err)
	}

	return id, hash, nil
}

func (s *Storage) CreateUserWithBalance(ctx context.Context, username, password string) (int, error) {
	var userID int

	sql := "SELECT create_user_with_balance($1, $2)"
	args := []interface{}{username, password}

	err := s.db.Pool.QueryRow(ctx, sql, args...).Scan(&userID)
	if err != nil {
		log.Errorf("repo.pgdb.CreateUserWithBalance.QueryRow: %v:", err)
		return 0, fmt.Errorf("could not create user: %w", err)
	}

	return userID, nil
}

func (s *Storage) GetPriceItem(ctx context.Context, item string) (int, error) {
	sql, args, _ := s.db.Builder.
		Select("price").
		From("items_price").
		Where("name = ?", item).
		ToSql()

	var price int
	err := s.db.Pool.QueryRow(ctx, sql, args...).Scan(&price)
	if err != nil {
		log.Errorf("repo.pgdb.GetPriceItem.QueryRow: %v:", err)
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, nil
		}
		return 0, err
	}

	return price, nil
}

func (s *Storage) UpdateBalance(ctx context.Context, userId int, coins int) error {
	sql, args, _ := s.db.Builder.
		Update("balances").
		Set("coins", squirrel.Expr("coins + ?", coins)).
		Where(squirrel.Eq{"user_id": userId}).
		ToSql()

	_, err := s.db.Pool.Exec(ctx, sql, args...)
	if err != nil {
		log.Errorf("repo.pgdb.UpdateBalance.Query: %v:", err)

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23514" {
			return fmt.Errorf("balance cannot be less than or equal to zero")
		}

		return err
	}
	return nil
}

func (s *Storage) AddItem(ctx context.Context, userId int, item string) error {
	sql, args, _ := s.db.Builder.
		Insert("user_inventory").
		Columns("user_id", "item_name", "quantity").
		Values(userId, item, 1).
		Suffix("ON CONFLICT (user_id, item_name) DO UPDATE " +
			"SET quantity = user_inventory.quantity + EXCLUDED.quantity").
		ToSql()

	_, err := s.db.Pool.Exec(ctx, sql, args...)
	if err != nil {
		log.Errorf("repo.pgdb.AddItem.Exec: %v:", err)
		return err
	}

	return nil
}

func (s *Storage) GetNameUser(ctx context.Context, userId int) (string, error) {
	sql, args, _ := s.db.Builder.
		Select("name").
		From("users").
		Where("id = ?", userId).
		ToSql()

	var name string
	err := s.db.Pool.QueryRow(ctx, sql, args...).Scan(&name)
	if err != nil {
		log.Errorf("repo.pgdb.GetNameUser.QueryRow: %v", err)
		if errors.Is(err, pgx.ErrNoRows) {
			return "", nil
		}
	}

	return name, nil
}

func (s *Storage) TransferCoins(ctx context.Context, fromUserId, toUserId, amount int) error {
	const fn = "repo.pgdb.TransferCoins"

	tx, err := s.db.Pool.Begin(ctx)
	if err != nil {
		log.Errorf("%s.Begin: %v", fn, err)
		return err
	}

	sql, args, _ := s.db.Builder.
		Update("balances").
		Set("coins", squirrel.Expr("CASE WHEN user_id = ? THEN coins - ? WHEN user_id = ? THEN coins + ? ELSE coins END", fromUserId, amount, toUserId, amount)).
		Where(squirrel.Eq{"user_id": []int{fromUserId, toUserId}}).
		ToSql()

	res, err := tx.Exec(ctx, sql, args...)
	if err != nil {
		log.Errorf("%s.balances.Exec: %v", fn, err)

		tx.Rollback(ctx)
		return err
	}

	rowAffected := res.RowsAffected()
	if rowAffected != 2 {
		log.Errorf("%s.RowsAffected: %v", fn, err)

		tx.Rollback(ctx)
		return fmt.Errorf("rows affected != 2")
	}

	sql, args, _ = s.db.Builder.
		Insert("transactions").
		Columns("sender_id", "receiver_id", "amount", "created_at").
		Values(fromUserId, toUserId, amount, time.Now()).
		ToSql()

	_, err = tx.Exec(ctx, sql, args...)
	if err != nil {
		log.Errorf("%s.transactions.Exec: %v", fn, err)

		tx.Rollback(ctx)
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		log.Errorf("%s.Commit: %v", fn, err)
		return err
	}

	return nil
}

func (s *Storage) GetReceivedTransactions(ctx context.Context, userId int) ([]repo.GetReceivedTransactionOutput, error) {
	sql, args, _ := s.db.Builder.
		Select("users.username AS from_user", "transactions.amount").
		From("transactions").
		Join("users ON transactions.sender_id = users.id").
		Where(squirrel.Eq{"transactions.receiver_id": userId}).
		ToSql()

	rows, err := s.db.Pool.Query(ctx, sql, args...)
	if err != nil {
		log.Errorf("repo.pgdb.GetReceivedTransactions.Query: %v", err)
		return nil, err
	}
	defer rows.Close()

	var transactions []repo.GetReceivedTransactionOutput
	for rows.Next() {
		var transaction repo.GetReceivedTransactionOutput
		if err = rows.Scan(&transaction.FromUser, &transaction.Amount); err != nil {
			log.Errorf("repo.pgdb.GetReceivedTransactions.Scan: %v", err)
			return nil, err
		}

		transactions = append(transactions, transaction)
	}

	return transactions, nil
}

func (s *Storage) GetSentTransactions(ctx context.Context, userId int) ([]repo.GetSentTransactionOutput, error) {
	sql, args, _ := s.db.Builder.
		Select("users.username AS to_user", "transactions.amount").
		From("transactions").
		Join("users ON transactions.receiver_id = users.id").
		Where(squirrel.Eq{"transactions.sender_id": userId}).
		ToSql()

	rows, err := s.db.Pool.Query(ctx, sql, args...)
	if err != nil {
		log.Errorf("repo.pgdb.GetSentTransactions.Query: %v", err)
		return nil, err
	}
	defer rows.Close()

	var transactions []repo.GetSentTransactionOutput
	for rows.Next() {
		var transaction repo.GetSentTransactionOutput
		if err = rows.Scan(&transaction.ToUser, &transaction.Amount); err != nil {
			log.Errorf("repo.pgdb.GetSentTransactions.Scan: %v", err)
			return nil, err
		}

		transactions = append(transactions, transaction)
	}

	return transactions, nil
}

func (s *Storage) GetBalance(ctx context.Context, userId int) (int, error) {
	sql, args, _ := s.db.Builder.
		Select("coins").
		From("balances").
		Where("user_id = ?", userId).
		ToSql()

	var coins int
	err := s.db.Pool.QueryRow(ctx, sql, args...).Scan(&coins)
	if err != nil {
		log.Errorf("repo.pgdb.GetBalance.QueryRow: %v", err)
		return 0, err
	}

	return coins, nil
}

func (s *Storage) GetItems(ctx context.Context, userId int) ([]repo.Item, error) {
	sql, args, _ := s.db.Builder.
		Select("item_name", "quantity").
		From("user_inventory").
		Where("user_id = ?", userId).
		ToSql()

	rows, err := s.db.Pool.Query(ctx, sql, args...)
	if err != nil {
		log.Errorf("repo.pgdb.GetItems.Query: %v", err)
		return nil, err
	}
	defer rows.Close()

	var items []repo.Item
	for rows.Next() {
		var item repo.Item
		if err = rows.Scan(&item.Type, &item.Quantity); err != nil {
			log.Errorf("repo.pgdb.GetItems.Scan: %v", err)
			return nil, err
		}

		items = append(items, item)
	}

	return items, nil
}
