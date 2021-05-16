package database

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/newrelic/go-agent/v3/integrations/nrpq"
)

type Conn interface {
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	Begin() (TxConn, error)
	Transaction(func(TxConn) error) error
	Close() error
}

type ConnPool struct {
	*sqlx.DB
}

func New(driver string, conn string) (*ConnPool, error) {
	db, err := sqlx.Open(driver, conn)
	if err != nil {
		return nil, err
	}

	db.SetConnMaxLifetime(30 * time.Second)
	db.SetMaxOpenConns(100)
	db.SetMaxIdleConns(10)

	return &ConnPool{db}, nil
}

func (c ConnPool) Transaction(f func(TxConn) error) error {
	tx, err := c.Begin()
	if err != nil {
		return err
	}

	if err = f(tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (c ConnPool) Begin() (TxConn, error) {
	tx, err := c.DB.Beginx()
	if err != nil {
		return nil, err
	}

	return &TxDatabase{c, tx}, nil
}
