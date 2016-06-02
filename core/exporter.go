package core

import (
	"../transaction"
	"database/sql"
	"fmt"
	"io"
	"os"
)

func safe_string(s string) string {
	if len(s) == 0 {
		s = "-"
	}
	return s
}

func ExportLog(db *sql.DB, w io.Writer) {
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
