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
	"fmt"
	"github.com/DataDrake/ApacheLog2DB/transaction"
	"github.com/jmoiron/sqlx"
	"io"
	"os"
)

func safe_string(s string) string {
	if len(s) == 0 {
		s = "-"
	}
	return s
}

func ExportLog(db *sqlx.DB, w io.Writer) {
	txns, err := transaction.ReadAll(db)
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(1)
	}
	for _, txn := range txns {
		request := ""
		lV := len(txn.Verb)
		lU := len(txn.Dest.URI)
		lP := len(txn.Protocol)
		if lV > 0 {
			request = txn.Verb

		}
		if len(request) > 0 && lU > 0 {
			request = request + " " + txn.Dest.URI
		} else if lU > 0 {
			request = txn.Dest.URI
		}

		if len(request) > 0 && lP > 0 {
			request = request + " " + txn.Protocol
		} else if lP > 0 {
			request = txn.Protocol
		}

		fmt.Fprintf(
			w,
			APACHE_COMBINED_FORMAT,
			safe_string(txn.Source.IP),
			safe_string(txn.Ident),
			safe_string(txn.User.Name),
			txn.Occurred.Format(APACHE_TIME_LAYOUT),
			request,
			txn.Status,
			txn.Size,
			txn.Referrer,
			safe_string(txn.Agent.Name),
		)
	}
}
