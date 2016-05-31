package core

import (
	"ApacheLog2DB/agent"
	"ApacheLog2DB/destination"
	"ApacheLog2DB/source"
	"ApacheLog2DB/transaction"
	"ApacheLog2DB/user"
	"bufio"
	"database/sql"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"
)

func ImportLog(log io.Reader, db *sql.DB) {
	err := CreateAllTables(db)
	if err != nil {
		fmt.Fprintf(os.Stderr,err.Error())
		os.Exit(1)
	}
	scan := bufio.NewScanner(log)
	for scan.Scan() {
		l := scan.Text()
		line := APACHE_COMBINED.FindStringSubmatch(l)
		if line != nil {

			//Get Source
			src, err := source.ReadOrCreate(db, line[1])
			if err != nil {
				fmt.Fprintf(os.Stderr, "[Source] Warning: %s\n", err.Error())
				continue
			}

			//Get ident
			ident := line[2]

			//Get User
			u, err := user.ReadOrCreate(db, line[3])
			if err != nil {
				fmt.Fprintf(os.Stderr, "[User] Warning: %s\n", err.Error())
				continue
			}

			//Get Time Occurred
			occurred, err := time.Parse(APACHE_TIME_LAYOUT, line[4])
			if err != nil {
				fmt.Fprintf(os.Stderr, "[Occured] Warning: %s\n", err.Error())
				continue
			}

			// Parse request string
			request := APACHE_REQUEST.FindStringSubmatch(line[5])
			verb := request[1]
			uri := request[2]
			protocol := request[3]


			dest, err := destination.ReadOrCreate(db, uri)
			if err != nil {
				fmt.Fprintf(os.Stderr, "[Destination] Warning: %s\n", err.Error())
				continue
			}

			//Convert status code to integer
			status, err := strconv.Atoi(line[6])
			if err != nil {
				fmt.Fprintf(os.Stderr, "[Status] Warning: %s\n", err.Error())
				continue
			}

			//Convert size to integer
			size, err := strconv.Atoi(line[7])
			if err != nil {
				fmt.Fprintf(os.Stderr, "[Size] Warning: %s\n", err.Error())
				continue
			}

			//Get Referred
			referrer := line[8]

			a, err := agent.ReadOrCreate(db, line[9])
			if err != nil {
				fmt.Fprintf(os.Stderr, "[Agent] Warning: %s\n", err.Error())
				continue
			}

			trans := &transaction.Transaction{-1, ident, verb, protocol, status, size, referrer, occurred, src, dest, a, u}

			err = transaction.Insert(db, trans)
			if err != nil {
				fmt.Fprintf(os.Stderr, "[TXN] Warning: %s\n", err.Error())
				continue
			}
		} else {
			fmt.Fprintf(os.Stderr,"Line did not match: %s\n", l)
		}
	}
}
