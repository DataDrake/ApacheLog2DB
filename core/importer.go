//
// Copyright 2016-2017 Bryan T. Meyers <bmeyers@datadrake.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package core

import (
	"bufio"
	"fmt"
	"github.com/DataDrake/ApacheLog2DB/agent"
	"github.com/DataDrake/ApacheLog2DB/destination"
	"github.com/DataDrake/ApacheLog2DB/source"
	"github.com/DataDrake/ApacheLog2DB/transaction"
	"github.com/DataDrake/ApacheLog2DB/user"
	"github.com/jmoiron/sqlx"
	"io"
	"os"
	"strconv"
	"time"
)

func ImportLog(log io.Reader, db *sqlx.DB) {
	err := CreateAllTables(db)
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(1)
	}
	scan := bufio.NewScanner(log)
	for scan.Scan() {
		l := scan.Text()
		line := APACHE_COMBINED.FindStringSubmatch(l)
		if line != nil {
			t := &transaction.Transaction{}
			//Get Source
			t.Source, err = source.ReadOrCreate(db, line[1])
			if err != nil {
				fmt.Fprintf(os.Stderr, "[Source] Warning: %s\n", err.Error())
				continue
			}

			//Get ident
			t.Ident = line[2]

			//Get User
			t.User, err = user.ReadOrCreate(db, line[3])
			if err != nil {
				fmt.Fprintf(os.Stderr, "[User] Warning: %s\n", err.Error())
				continue
			}

			//Get Time Occurred
			t.Occurred, err = time.Parse(APACHE_TIME_LAYOUT, line[4])
			if err != nil {
				fmt.Fprintf(os.Stderr, "[Occured] Warning: %s\n", err.Error())
				continue
			}

			// Parse request string
			request := APACHE_REQUEST.FindStringSubmatch(line[5])
			t.Verb = request[1]
			uri := request[2]
			t.Protocol = request[3]

			t.Dest, err = destination.ReadOrCreate(db, uri)
			if err != nil {
				fmt.Fprintf(os.Stderr, "[Destination] Warning: %s\n", err.Error())
				continue
			}

			//Convert status code to integer
			t.Status, err = strconv.Atoi(line[6])
			if err != nil {
				fmt.Fprintf(os.Stderr, "[Status] Warning: %s\n", err.Error())
				continue
			}

			//Convert size to integer
			t.Size, err = strconv.Atoi(line[7])
			if err != nil {
				fmt.Fprintf(os.Stderr, "[Size] Warning: %s\n", err.Error())
				continue
			}

			//Get Referred
			t.Referrer = line[8]

			t.Agent, err = agent.ReadOrCreate(db, line[9])
			if err != nil {
				fmt.Fprintf(os.Stderr, "[Agent] Warning: %s\n", err.Error())
				continue
			}

			err = transaction.Insert(db, t)
			if err != nil {
				fmt.Fprintf(os.Stderr, "[TXN] Warning: %s\n", err.Error())
				continue
			}
		} else {
			fmt.Fprintf(os.Stderr, "Line did not match: %s\n", l)
		}
	}
}
