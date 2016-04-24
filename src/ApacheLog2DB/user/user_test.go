package user

import (
	"testing"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
)

func TestTableCreate(t *testing.T){
	os.Remove("/tmp/foo.db")
	db, err := sql.Open("sqlite3", "/tmp/foo.db")
	if err != nil {
		log.Fatal(err)
	}
	err = CreateTable(db)
	if err != nil {t.Error(err.Error())}
	defer db.Close()
}

func TestReadAllEmpty(t *testing.T){
	os.Remove("/tmp/foo.db")
	db, err := sql.Open("sqlite3", "/tmp/foo.db")
	if err != nil {
		log.Fatal(err)
	}
	err = CreateTable(db)
	if err != nil {t.Error(err.Error())}
	sources,err := ReadAll(db)
	if err != nil {t.Error(err.Error())}
	if len(sources) > 0 {t.Error("Table should be empty")}
	defer db.Close()
}

func TestInsertOne(t *testing.T) {
	os.Remove("/tmp/foo.db")
	db, err := sql.Open("sqlite3", "/tmp/foo.db")
	if err != nil {
		log.Fatal(err)
	}
	err = CreateTable(db)
	if err != nil {t.Error(err.Error())}
	sources,err := ReadAll(db)
	if err != nil {t.Error(err.Error())}
	if len(sources) > 0 {t.Error("Table should be empty")}

	err = Insert(db,&User{0,"abc1234"})
	if err != nil {t.Error(err.Error())}

	sources,err = ReadAll(db)
	if err != nil {t.Error(err.Error())}
	if len(sources) != 1 {t.Errorf("Table should have 1 entry, found: %d", len(sources))}
	defer db.Close()
}
