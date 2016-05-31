package ApacheLog2DB

import (
	"ApacheLog2DB/agent"
	"ApacheLog2DB/destination"
	"ApacheLog2DB/source"
	"ApacheLog2DB/transaction"
	"ApacheLog2DB/user"
	"database/sql"
	"errors"
	_ "github.com/mattn/go-sqlite3"
	"regexp"
	"io"
"strings"
	"time"
	"strconv"
)

var LOG2DB_TABLES = []string{"destinations", "sources", "transactions", "users", "user_agents"}
var APACHE_COMBINED = regexp.MustCompile("^(\\S*).(\\S*).(\\S*).\\[(.*)\\].\"([^\"]*)\".(\\d{3}).(\\d*).\"([^\"]*)\".\"([^\"]*)\"$")
var APACHE_TIME_LAYOUT = "2006-01-02T15:04:05.000Z"

func SliceContains(vs []string, v string) bool {
	for _, curr := range vs {
		if curr == v {
			return true
		}
	}
	return false
}

func get_tables(tx *sql.Tx) ([]string, error) {
	tables := make([]string, 0)
	found, err := tx.Query("SELECT name FROM sqlite_master WHERE type='table'")
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
	err, missing := CheckTables(db)

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

	if SliceContains(missing, "transactions") {
		err = transaction.CreateTable(db)
	}
	return err
}

func CheckTables(db *sql.DB) (error, []string) {

	tx, err := db.Begin()
	if err != nil {
		return err, nil
	}
	tables, err := get_tables(tx)
	tx.Commit()

	missing := make([]string, 0)
	for _, t := range LOG2DB_TABLES {
		if !SliceContains(tables, t) {
			missing = append(missing, t)
		}
	}
	if len(missing) != 0 {
		return errors.New("Missing tables"), missing
	}
	return nil, missing
}

func ImportLog(log string) error{
	for _,line := range APACHE_COMBINED.FindAllStringSubmatch(log,-1) {
		source := line[1]
		ident := line[2]
		username := line[3]
		occurred,err := time.Parse(APACHE_TIME_LAYOUT,line[4])
		if err != nil {return err}
		request := strings.Split(line[5], " ")

		var verb string
		var uri string
		var protocol string
		switch len(request) {
		case 3:
			verb = request[0]
			uri = request[1]
			protocol = request[2]
		case 2:
			uri = request[0]
			protocol = request[1]
		case 1:
			uri = request[0]
		}

		status,err := strconv.Atoi(line[6])
		if err != nil {return err}

		size,err := strconv.Atoi(line[7])
		if err != nil {return err}

		referrer := line[8]

		agentname := line[9]

	}
}

func ExportLog(w *io.Writer) {

}
