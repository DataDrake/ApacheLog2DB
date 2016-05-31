package main

import (
	"ApacheLog2DB/core"
	"database/sql"
	"flag"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"io"
	"os"
)

func usage() {
	fmt.Println("USAGE: log2db [OPTION]... SOURCE DEST")
	flag.PrintDefaults()
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
	var reader *io.Reader
	var writer *io.Writer

	if !export {
		if (args[0] == "-") || (args[0] == "--") {
			reader = os.Stdin
		} else {
			reader, err = os.Open(args[0])
		}
	} else {
		if !(args[0] == "-") && !(args[0] == "--") {
			fmt.Fprintf(os.Stderr, "Input file must be a db string")
			os.Exit(1)
		}
		db, err = sql.Open("sqlite3", args[0])
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(1)
	}

	if !export {
		if (args[1] == "-") || (args[1] == "--") {
			fmt.Fprintf(os.Stderr, "Output file must be a db string")
			os.Exit(1)
		}
		db, err = sql.Open("sqlite3", args[1])
	} else {
		if (args[1] == "-") || (args[1] == "--") {
			writer = os.Stdout
		} else {
			writer, err = os.Open(args[1])
		}
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(1)
	}

	if export {
		core.ExportLog(db, writer)
	} else {
		core.ImportLog(reader, db)
	}

	os.Exit(0)
}
