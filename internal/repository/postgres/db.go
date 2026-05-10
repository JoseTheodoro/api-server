package postgres

import (
	"context"
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type Connection struct {
	dsn string
}

func NewConnection(dsn string) *Connection {
	return &Connection{dsn: dsn}
}

func (c *Connection) Connect(ctx context.Context) (*sql.DB, error) {

	db, err := sql.Open("pgx", c.dsn)
	if err != nil {
		return nil, err
	}
	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, err
	}
	return db, nil
}
