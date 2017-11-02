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
	"errors"
	"github.com/DataDrake/ApacheLog2DB/agent"
	"github.com/DataDrake/ApacheLog2DB/destination"
	"github.com/DataDrake/ApacheLog2DB/global"
	"github.com/DataDrake/ApacheLog2DB/source"
	"github.com/DataDrake/ApacheLog2DB/transaction"
	"github.com/DataDrake/ApacheLog2DB/user"
	"github.com/jmoiron/sqlx"
)

var GetTablesQueries = map[string]string{
	"mysql":  "SHOW TABLES",
	"sqlite": "SELECT name FROM sqlite_master WHERE type='table'",
}

func GetTables(db *sqlx.DB) ([]string, error) {
	tables := make([]string, 0)
	found, err := db.Query(GetTablesQueries[global.DB_TYPE])
	if err != nil {
		return nil, err
	}
	var table string
	found.Scan(&table)
	if len(table) > 0 {
		tables = append(tables, table)
	}
	for found.Next() {
		found.Scan(&table)
		if len(table) > 0 {
			tables = append(tables, table)
		}
	}
	found.Close()
	return tables, err
}

func CreateAllTables(db *sqlx.DB) error {
	missing, err := CheckTables(db)

	if SliceContains(missing, "user_agents") {
		err = agent.CreateTable(db)
		if err != nil {
			return err
		}
	}

	if SliceContains(missing, "sources") {
		err = source.CreateTable(db)
		if err != nil {
			return err
		}
	}

	if SliceContains(missing, "destinations") {
		err = destination.CreateTable(db)
		if err != nil {
			return err
		}
	}

	if SliceContains(missing, "users") {
		err = user.CreateTable(db)
		if err != nil {
			return err
		}
	}

	if SliceContains(missing, "txns") {
		err = transaction.CreateTable(db)
	}
	return err
}

func CheckTables(db *sqlx.DB) ([]string, error) {
	tables, err := GetTables(db)
	if err != nil {
		return nil, err
	}
	missing := make([]string, 0)
	for _, t := range LOG2DB_TABLES {
		if !SliceContains(tables, t) {
			missing = append(missing, t)
		}
	}
	if len(missing) != 0 {
		return missing, errors.New("Missing tables")
	}
	return missing, nil
}
