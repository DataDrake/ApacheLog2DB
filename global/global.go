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

package global

import (
	"errors"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"strings"
)

var DB_TYPE string

func OpenDatabase(connection string) (db *sqlx.DB, err error) {
	protocol := strings.Split(connection, ":")
	switch protocol[0] {
	case "mysql":
		DB_TYPE = "mysql"
		db, err = sqlx.Open("mysql", strings.Replace(connection, "mysql://", "", 1))
	case "sqlite":
		DB_TYPE = "sqlite"
		db, err = sqlx.Open("sqlite3", strings.Replace(protocol[1], "//", "", 1))
	default:
		err = errors.New("Unsupported database type: " + protocol[0])
	}
	return
}
