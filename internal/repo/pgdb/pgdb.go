package pgdb

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/magmaheat/merchStore/pkg/postgres"
	log "github.com/sirupsen/logrus"
)

type Storage struct {
	db *postgres.Postgres
}

func New(db *postgres.Postgres) *Storage {
	return &Storage{db: db}
}

func (s *Storage) CreateUser(ctx context.Context, username string, password string) (int, error) {
	sql, args, _ := s.db.Builder.
		Insert("users").
		Columns("username, password").
		Values(username, password).
		ToSql()

	var id int
	err := s.db.Pool.QueryRow(ctx, sql, args...).Scan(&id)
	if err != nil {
		log.Errorf("repo.CreateUser.QueryRow: %v:", err)
		return 0, fmt.Errorf("could not create user: %w", err)
	}

	return id, nil
}

func (s *Storage) GetUser(ctx context.Context, username string) (int, string, error) {
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
		log.Errorf("repo.GetUser.QueryRow: %v:", err)
		return 0, "", fmt.Errorf("could not get user: %w", err)
	}

	return id, hash, nil
}

func (s *Storage) CreateUserWithBalance(ctx context.Context, username string, password string) (int, error) {
	tx, err := s.db.Pool.Begin(ctx)
	if err != nil {
		return 0, err
	}

	sql, args, _ := s.db.Builder.
		Insert("users").
		Columns("username, password").
		Values(username, password).
		Suffix("RETURNING id").
		ToSql()

	var id int
	err = tx.QueryRow(ctx, sql, args...).Scan(&id)
	if err != nil {
		log.Errorf("repo.CreateUserWithBalance.QueryRow: %v:", err)
		tx.Rollback(ctx)
		return 0, fmt.Errorf("could not create user: %w", err)
	}

	sql, args, _ = s.db.Builder.
		Insert("balances").
		Columns("user_id").
		Values(id).
		ToSql()

	_, err = tx.Exec(ctx, sql, args...)
	if err != nil {
		log.Errorf("repo.CreateUserWithBalance.Exec: %v:", err)
		tx.Rollback(ctx)
		return 0, fmt.Errorf("could not create balance: %w", err)
	}

	if err = tx.Commit(ctx); err != nil {
		return 0, err
	}

	return id, nil
}
