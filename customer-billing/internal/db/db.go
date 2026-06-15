package db

import (
	"context"
	"database/sql"
	"embed"
	"fmt"

	sqlc "github.com/andrewrutherfoord/fed-bill-poc/customer-billing/internal/db/sqlc"
	_ "github.com/glebarez/go-sqlite"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

type DB struct {
	*sqlc.Queries
	conn *sql.DB
}

func Open(path string) (*DB, error) {
	conn, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("opening database: %w", err)
	}

	if err := conn.Ping(); err != nil {
		return nil, fmt.Errorf("pinging database: %w", err)
	}

	return &DB{
		Queries: sqlc.New(conn),
		conn:    conn,
	}, nil
}

func (d *DB) Close() error {
	return d.conn.Close()
}

func (d *DB) WithTx(ctx context.Context, fn func(*sqlc.Queries) error) error {
	tx, err := d.conn.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	if err := fn(sqlc.New(tx)); err != nil {
		_ = tx.Rollback()
		return err
	}
	return tx.Commit()
}
