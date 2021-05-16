package database

import (
	"github.com/spf13/viper"
)

var conn Conn

func InitPool() error {
	defaultConn, err := New(
		viper.GetString("DATABASE_DRIVER"),
		viper.GetString("DATABASE_DSN"),
	)
	if err != nil {
		return err
	}

	conn = defaultConn

	return nil
}

func Pool() Conn {
	return conn
}

func ClosePool() error {
	return conn.Close()
}
