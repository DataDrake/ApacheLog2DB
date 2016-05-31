package core

import (
	"bufio"
	"database/sql"
	"io"
	"strconv"
	"strings"
	"time"
	"fmt"
	"os"
)

func ImportLog(log *io.Reader, db *sql.DB) error {
	scan := bufio.NewScanner(*log)
	for scan.Scan() {
		line := APACHE_COMBINED.FindStringSubmatch(scan.Text())
		if line != nil {
			source := line[1]
			ident := line[2]
			username := line[3]
			occurred, err := time.Parse(APACHE_TIME_LAYOUT, line[4])
			if err != nil {
				fmt.Fprintf(os.Stderr,"Warning: %s\n",err.Error())
				continue
			}
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

			status, err := strconv.Atoi(line[6])
			if err != nil {
				fmt.Fprintf(os.Stderr,"Warning: %s\n",err.Error())
				continue
			}

			size, err := strconv.Atoi(line[7])
			if err != nil {
				fmt.Fprintf(os.Stderr,"Warning: %s\n",err.Error())
				continue
			}

			referrer := line[8]

			agentname := line[9]

		}
	}
	return nil
}
