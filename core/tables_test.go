package core

import (
	"database/sql"
	"github.com/DataDrake/ApacheLog2DB/agent"
	"github.com/DataDrake/ApacheLog2DB/destination"
	"github.com/DataDrake/ApacheLog2DB/source"
	"github.com/DataDrake/ApacheLog2DB/transaction"
	"github.com/DataDrake/ApacheLog2DB/user"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"testing"
)

func TestCheckTables(t *testing.T) {
	os.Remove("/tmp/foo.db")

	db, err := sql.Open("sqlite3", "/tmp/foo.db")
	defer db.Close()

	if err != nil {
		t.Error(err.Error())
	}

	missing, err := CheckTables(db)
	if err == nil {
		t.Error("Should not have succeeded")
	}
	if len(missing) != len(LOG2DB_TABLES) {
		t.Error("Should have found all tables missing.")
	}

	err = user.CreateTable(db)
	if err != nil {
		t.Error(err.Error())
	}
	missing, err = CheckTables(db)
	if err == nil {
		t.Error("Should not have succeeded")
	}
	if len(missing) != (len(LOG2DB_TABLES) - 1) {
		t.Error("Should have found all but one tables missing.")
	}
}

func TestCheckTablesComplete(t *testing.T) {
	os.Remove("/tmp/foo.db")

	db, err := sql.Open("sqlite3", "/tmp/foo.db")
	defer db.Close()

	if err != nil {
		t.Error(err.Error())
	}

	missing, err := CheckTables(db)
	if err == nil {
		t.Error("Should not have succeeded")
	}
	if len(missing) != len(LOG2DB_TABLES) {
		t.Error("Should have found all tables missing.")
	}

	err = agent.CreateTable(db)
	if err != nil {
		t.Error(err.Error())
	}

	err = destination.CreateTable(db)
	if err != nil {
		t.Error(err.Error())
	}

	err = source.CreateTable(db)
	if err != nil {
		t.Error(err.Error())
	}

	err = user.CreateTable(db)
	if err != nil {
		t.Error(err.Error())
	}

	err = transaction.CreateTable(db)
	if err != nil {
		t.Error(err.Error())
	}

	missing, err = CheckTables(db)
	if err != nil {
		t.Error("Should have succeeded")
	}
	if len(missing) != 0 {
		t.Error("Should have found no tables missing.")
	}
}

func TestCreateAllTables(t *testing.T) {
	os.Remove("/tmp/foo.db")

	db, err := sql.Open("sqlite3", "/tmp/foo.db")
	defer db.Close()

	if err != nil {
		t.Error(err.Error())
	}

	missing, err := CheckTables(db)
	if err == nil {
		t.Error("Should not have succeeded")
	}
	if len(missing) != len(LOG2DB_TABLES) {
		t.Error("Should have found all tables missing.")
	}

	err = CreateAllTables(db)

	missing, err = CheckTables(db)
	if err != nil {
		t.Error("Should have succeeded")
	}
	if len(missing) != 0 {
		t.Error("Should have found no tables missing.")
	}
}
