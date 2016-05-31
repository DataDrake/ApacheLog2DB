package core

import (
	"database/sql"
	"io"
	"ApacheLog2DB/transaction"
	"fmt"
	"os"
)

func safe_string(s string) string {
	if len(s) == 0 {
		s = "-"
	}
	return s
}

func ExportLog(db *sql.DB, w io.Writer) {
	txns,err := transaction.ReadAll(db)
	if err != nil {
		fmt.Fprintf(os.Stderr,err.Error())
		os.Exit(1)
	}
	for _,txn := range txns {
		fmt.Fprintln(
			w,
			APACHE_COMBINED_FORMAT,
			safe_string(txn.Source.IP),
			safe_string(txn.Ident),
			safe_string(txn.User.Name),
			txn.Occurred.Format(APACHE_TIME_LAYOUT),
			txn.Verb,
			txn.Dest.URI,
			txn.Protocol,
			txn.Status,
			txn.Size,
			safe_string(txn.Referrer),
			safe_string(txn.Agent.Name),
		)
	}
}
