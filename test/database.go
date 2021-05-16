package test

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	"github.com/Masterminds/squirrel"
	"github.com/diegodesousas/clean-boilerplate-go/infra/database"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

func PrepareDatabaseTest() error {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("DATABASE_TEST_USER"),
		os.Getenv("DATABASE_TEST_PASSWORD"),
		os.Getenv("DATABASE_TEST_HOST"),
		os.Getenv("DATABASE_TEST_PORT"),
		os.Getenv("DATABASE_TEST_NAME"),
	)

	viper.Set("DATABASE_DRIVER", "postgres")
	viper.Set("DATABASE_DSN", dsn)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return errors.Wrapf(err, "error when trying connect to database test")
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return errors.Wrapf(err, "error when trying connect database test")
	}

	mg, err := migrate.NewWithDatabaseInstance(fmt.Sprintf("file://%s/infra/database/migrations", os.Getenv("BASE_PROJECT_PATH")), os.Getenv("DATABASE_TEST_NAME"), driver)
	if err != nil {
		return errors.Wrapf(err, "error when trying run migrations")
	}

	if err = mg.Up(); err != nil && err != migrate.ErrNoChange {
		return errors.Wrapf(err, "error when trying run migrations")
	}

	if err := database.InitPool(); err != nil {
		return errors.Wrapf(err, "error when trying to connect with the database")
	}

	return nil
}

func AutoRollback(f func(conn database.TxConn) error) error {
	txConn, err := database.Pool().Begin()
	if err != nil {
		return err
	}
	defer txConn.Rollback()

	if err := f(txConn); err != nil {
		return err
	}

	return nil
}

func RunFixtures(conn database.TxConn, tableName string, data ...map[string]interface{}) error {
	for _, clause := range data {
		query, args, err := squirrel.
			Insert(tableName).
			SetMap(clause).
			Suffix("RETURNING id").
			PlaceholderFormat(squirrel.Dollar).
			ToSql()

		if err != nil {
			return err
		}

		var test string
		err = conn.ExecContext(context.Background(), query, args...).Scan(&test)
		if err != nil {
			return err
		}
	}

	return nil
}
