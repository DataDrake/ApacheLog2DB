package global

import (
	"database/sql"
	"errors"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
	"strings"
)

var DB_TYPE string

func OpenDatabase(connection string) (db *sql.DB, err error) {
	protocol := strings.Split(connection, ":")
	switch protocol[0] {
	case "mysql":
		DB_TYPE = "mysql"
		db, err = sql.Open("mysql", strings.Replace(connection, "mysql://", "", 1))
	case "sqlite":
		DB_TYPE = "sqlite"
		db, err = sql.Open("sqlite3", strings.Replace(protocol[1], "//", "", 1))
	default:
		err = errors.New("Unsupported database type: " + protocol[0])
	}
	return
}
