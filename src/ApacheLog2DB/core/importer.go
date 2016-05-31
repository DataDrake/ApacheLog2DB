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
	"strings"
	"time"
)

func ImportLog(log *io.Reader, db *sql.DB) {
	scan := bufio.NewScanner(*log)
	for scan.Scan() {
		line := APACHE_COMBINED.FindStringSubmatch(scan.Text())
		if line != nil {

			//Get Source
			src, err := source.ReadOrCreate(db, line[1])
			if err != nil {
				fmt.Fprintf(os.Stderr, "Warning: %s\n", err.Error())
				continue
			}

			//Get ident
			ident := line[2]

			//Get User
			u, err := user.ReadOrCreate(db, line[3])
			if err != nil {
				fmt.Fprintf(os.Stderr, "Warning: %s\n", err.Error())
				continue
			}

			//Get Time Occurred
			occurred, err := time.Parse(APACHE_TIME_LAYOUT, line[4])
			if err != nil {
				fmt.Fprintf(os.Stderr, "Warning: %s\n", err.Error())
				continue
			}

			// Parse request string
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

			dest, err := destination.ReadOrCreate(db, uri)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Warning: %s\n", err.Error())
				continue
			}

			//Convert status code to integer
			status, err := strconv.Atoi(line[6])
			if err != nil {
				fmt.Fprintf(os.Stderr, "Warning: %s\n", err.Error())
				continue
			}

			//Convert size to integer
			size, err := strconv.Atoi(line[7])
			if err != nil {
				fmt.Fprintf(os.Stderr, "Warning: %s\n", err.Error())
				continue
			}

			//Get Referred
			referrer := line[8]

			a, err := agent.ReadOrCreate(db, line[9])
			if err != nil {
				fmt.Fprintf(os.Stderr, "Warning: %s\n", err.Error())
				continue
			}

			trans := &transaction.Transaction{-1, ident, verb, protocol, status, size, referrer, occurred, src, dest, a, u}

			err = transaction.Insert(db, trans)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Warning: %s\n", err.Error())
				continue
			}
		}
	}
}
