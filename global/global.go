package global

import (
	"errors"
	_ "github.com/go-sql-driver/mysql"
    "github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"strings"
)

var DB_TYPE string

func OpenDatabase(connection string) (db *sqlx.DB, err error) {
	protocol := strings.Split(connection, ":")
	switch protocol[0] {
	case "mysql":
		DB_TYPE = "mysql"
		db, err = sqlx.Open("mysql", strings.Replace(connection, "mysql://", "", 1))
	case "sqlite":
		DB_TYPE = "sqlite"
		db, err = sqlx.Open("sqlite3", strings.Replace(protocol[1], "//", "", 1))
	default:
		err = errors.New("Unsupported database type: " + protocol[0])
	}
	return
}
