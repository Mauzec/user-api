package db

import "github.com/jackc/pgx/v5/pgxpool"

// Store defines all methods to exec queries and transactions
// TODO: for future transactions
type Store interface {
	Querier
}

type PSQLSTore struct {
	db *pgxpool.Pool
	*Queries
}

func NewStore(db *pgxpool.Pool) Store {
	return &PSQLSTore{
		db:      db,
		Queries: New(db),
	}
}
