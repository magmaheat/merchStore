package pgdb

import (
	"context"
	"errors"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/magmaheat/merchStore/pkg/postgres"
	log "github.com/sirupsen/logrus"
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
		log.Errorf("repo.GetUserIdWithPassword.QueryRow: %v:", err)
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
		log.Errorf("repo.CreateUserWithBalance.QueryRow: %v:", err)
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
		log.Errorf("repo.GetPriceItem.QueryRow: %v:", err)
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, nil
		}
		return 0, err
	}

	return price, nil
}

func (s *Storage) UpdateBalance(ctx context.Context, userId int, price int) error {
	sql, args, _ := s.db.Builder.
		Update("balances").
		Set("money", squirrel.Expr("money + ?", price)).
		Where(squirrel.Eq{"user_id": userId}).
		ToSql()

	_, err := s.db.Pool.Exec(ctx, sql, args...)
	if err != nil {
		log.Errorf("repo.UpdateBalance.Query: %v:", err)

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
		ToSql()

	_, err := s.db.Pool.Exec(ctx, sql, args...)
	if err != nil {
		log.Errorf("repo.AddItem.Exec: %v:", err)
		return err
	}

	return nil
}
