package ApacheLog2DB
import (
	_ "github.com/mattn/go-sqlite3"
	"database/sql"
	"errors"
	"ApacheLog2DB/agent"
	"ApacheLog2DB/destination"
	"ApacheLog2DB/source"
"ApacheLog2DB/user"
	"ApacheLog2DB/transaction"
)

var LOG2DB_TABLES = []string{"destinations","sources","transactions","users","user_agents"}

func SliceContains(vs []string, v string) bool {
	for _,curr := range vs {
		if curr == v {return true}
	}
	return false
}

func get_tables(tx *sql.Tx) ([]string,error) {
	tables := make([]string,0)
	found,err := tx.Query("SELECT name FROM sqlite_master WHERE type='table'")
	if err != nil {
		return nil, err
	}
	var table string
	found.Scan(&table)
	if len(table) > 0 {
		tables = append(tables,table)
	}
	for found.Next() {
		found.Scan(&table)
		if len(table) > 0 {
			tables = append(tables,table)
		}
	}
	found.Close()
	return tables,err
}

func CreateAllTables(db *sql.DB) error {
	err,missing := CheckTables(db)

	if SliceContains(missing,"user_agents"){
		err = agent.CreateTable(db)
		if err != nil {return err}
	}

	if SliceContains(missing,"sources"){
		err = source.CreateTable(db)
		if err != nil {return err}
	}

	if SliceContains(missing,"destinations"){
		err = destination.CreateTable(db)
		if err != nil {return err}
	}

	if SliceContains(missing,"users"){
		err = user.CreateTable(db)
		if err != nil {return err}
	}

	if SliceContains(missing,"transactions"){
		err = transaction.CreateTable(db)
	}
	return err
}

func CheckTables(db *sql.DB) (error,[]string){

	tx,err := db.Begin()
	if err != nil {
		return err, nil
	}
	tables,err := get_tables(tx)
	tx.Commit()

	missing := make([]string,0)
	for _,t := range LOG2DB_TABLES {
		if !SliceContains(tables,t) {missing = append(missing,t)}
	}
	if len(missing) != 0 {
		return errors.New("Missing tables"),missing
	}
	return nil,missing
}