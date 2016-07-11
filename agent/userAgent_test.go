package agent

import (
	"github.com/DataDrake/ApacheLog2DB/global"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
	"testing"
)

func TestTableCreate(t *testing.T) {
	os.Remove("/tmp/foo.db")
	db, err := global.OpenDatabase("sqlite:///tmp/foo.db")
	if err != nil {
		log.Fatal(err)
	}
	err = CreateTable(db)
	if err != nil {
		t.Error(err.Error())
	}
	defer db.Close()
}

func TestReadAllEmpty(t *testing.T) {
	os.Remove("/tmp/foo.db")
	db, err := global.OpenDatabase("sqlite:///tmp/foo.db")
	if err != nil {
		log.Fatal(err)
	}
	err = CreateTable(db)
	if err != nil {
		t.Error(err.Error())
	}
	agents, err := ReadAll(db)
	if err != nil {
		t.Error(err.Error())
	}
	if len(agents) > 0 {
		t.Error("Table should be empty")
	}
	defer db.Close()
}

func TestInsertOne(t *testing.T) {
	os.Remove("/tmp/foo.db")
	db, err := global.OpenDatabase("sqlite:///tmp/foo.db")
	if err != nil {
		log.Fatal(err)
	}
	err = CreateTable(db)
	if err != nil {
		t.Error(err.Error())
	}
	agents, err := ReadAll(db)
	if err != nil {
		t.Error(err.Error())
	}
	if len(agents) > 0 {
		t.Error("Table should be empty")
	}

	err = Insert(db, &UserAgent{0, "Firefox"})
	if err != nil {
		t.Error(err.Error())
	}

	agents, err = ReadAll(db)
	if err != nil {
		t.Error(err.Error())
	}
	if len(agents) != 1 {
		t.Errorf("Table should have 1 entry, found: %d", len(agents))
	}
	defer db.Close()
}
