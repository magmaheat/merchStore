package pgdb

import "github.com/magmaheat/merchStore/pkg/postgres"

type Storage struct {
	db *postgres.Postgres
}

func New(db *postgres.Postgres) *Storage {
	return &Storage{db: db}
}
