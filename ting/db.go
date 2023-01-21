package ting

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

const driverName = "mysql"

var db *sql.DB

func InitDb(config DbConfig) error {
	var err error
	db, err = sql.Open(driverName, fmt.Sprintf("%s:%s@tcp(%s:%d)/ting?parseTime=true", config.UserName, config.Password, config.Host, config.Port))

	if err != nil {
		return err
	}

	err = db.Ping()

	if err != nil {
		return err
	}

	return nil
}

func CloseDb() error {
	if db == nil {
		return errors.New("database is not initialized")
	}

	return db.Close()
}

func Prepare(sql string) (*sql.Stmt, error) {
	if db == nil {
		return nil, errors.New("database is not initialized")
	}

	return db.Prepare(sql)
}
