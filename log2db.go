package main

import (
	"database/sql"
	"flag"
	"fmt"
	"github.com/DataDrake/ApacheLog2DB/core"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"strings"
)

func usage() {
	fmt.Println("USAGE: log2db [OPTION]... SOURCE DEST")
	flag.PrintDefaults()
}

func dbconnection(conn string) (*sql.DB, error) {
	if strings.HasPrefix(conn, "mysql://") {
		conn = strings.Replace(conn, "mysql://", "", 1)
		return sql.Open("mysql", conn)
	}
	return sql.Open("sqlite3", conn)
}

func main() {
	flag.Usage = func() { usage() }
	export := flag.Bool("e", false, "Export from SOURCE to DEST")
	flag.Parse()

	args := flag.Args()

	if len(args) != 2 {
		usage()
		os.Exit(1)
	}

	var db *sql.DB
	var err error
	reader := os.Stdin
	writer := os.Stdout

	if !*export {
		if !(args[0] == "-" || args[0] == "--") {
			reader, err = os.Open(args[0])
			defer reader.Close()
		}
	} else {
		if args[0] == "-" || args[0] == "--" {
			fmt.Fprintln(os.Stderr, "Input file must be a db string")
			os.Exit(1)
		}
		db, err = dbconnection(args[0])
		defer db.Close()
	}

	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	if !*export {
		if args[1] == "-" || args[1] == "--" {
			fmt.Fprintln(os.Stderr, "Output file must be a db string")
			os.Exit(1)
		}
		db, err = dbconnection(args[1])
		defer db.Close()
	} else {
		if !(args[1] == "-" || args[1] == "--") {
			writer, err = os.Create(args[1])
			defer writer.Close()
		}
	}

	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	if *export {
		core.ExportLog(db, writer)
	} else {
		core.ImportLog(reader, db)
	}
	os.Exit(0)
}
