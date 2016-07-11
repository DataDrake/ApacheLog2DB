package core

import (
	"database/sql"
	"errors"
	"github.com/DataDrake/ApacheLog2DB/agent"
	"github.com/DataDrake/ApacheLog2DB/destination"
	"github.com/DataDrake/ApacheLog2DB/source"
	"github.com/DataDrake/ApacheLog2DB/transaction"
	"github.com/DataDrake/ApacheLog2DB/user"
)

var DB_TYPE string

var GET_TABLES = map[string]string{
	"mysql":"SHOW TABLES",
	"sqlite":"SELECT name FROM sqlite_master WHERE type='table' ORDER BY name;",
}

func get_tables(db *sql.DB) ([]string, error) {
	tables := make([]string, 0)
	//found, err := db.Query("SELECT name FROM sqlite_master WHERE type='table'")
	found, err := db.Query(DB_TYPE)
	if err != nil {
		return nil, err
	}
	var table string
	found.Scan(&table)
	if len(table) > 0 {
		tables = append(tables, table)
	}
	for found.Next() {
		found.Scan(&table)
		if len(table) > 0 {
			tables = append(tables, table)
		}
	}
	found.Close()
	return tables, err
}

func CreateAllTables(db *sql.DB) error {
	missing, err := CheckTables(db)

	if SliceContains(missing, "user_agents") {
		err = agent.CreateTable(db)
		if err != nil {
			return err
		}
	}

	if SliceContains(missing, "sources") {
		err = source.CreateTable(db)
		if err != nil {
			return err
		}
	}

	if SliceContains(missing, "destinations") {
		err = destination.CreateTable(db)
		if err != nil {
			return err
		}
	}

	if SliceContains(missing, "users") {
		err = user.CreateTable(db)
		if err != nil {
			return err
		}
	}

	if SliceContains(missing, "txns") {
		err = transaction.CreateTable(db)
	}
	return err
}

func CheckTables(db *sql.DB) ([]string, error) {
	tables, err := get_tables(db)
	if err != nil {
		return nil, err
	}
	missing := make([]string, 0)
	for _, t := range LOG2DB_TABLES {
		if !SliceContains(tables, t) {
			missing = append(missing, t)
		}
	}
	if len(missing) != 0 {
		return missing, errors.New("Missing tables")
	}
	return missing, nil
}
