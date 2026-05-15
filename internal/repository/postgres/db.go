package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Connection struct {
	dsn string
}

func NewConnection(dsn string) *Connection {
	return &Connection{dsn: dsn}
}

func (c *Connection) Connect(ctx context.Context) (*pgxpool.Pool, error) {

	db, err := pgxpool.New(ctx, c.dsn)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(ctx); err != nil {
		db.Close()
		return nil, err
	}
	return db, nil
}
