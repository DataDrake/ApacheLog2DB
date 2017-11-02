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
	"github.com/DataDrake/ApacheLog2DB/agent"
	"github.com/DataDrake/ApacheLog2DB/destination"
	"github.com/DataDrake/ApacheLog2DB/global"
	"github.com/DataDrake/ApacheLog2DB/source"
	"github.com/DataDrake/ApacheLog2DB/transaction"
	"github.com/DataDrake/ApacheLog2DB/user"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"testing"
)

func TestCheckTables(t *testing.T) {
	os.Remove("/tmp/foo.db")

	db, err := global.OpenDatabase("sqlite:///tmp/foo.db")
	defer db.Close()

	if err != nil {
		t.Error(err.Error())
	}

	missing, err := CheckTables(db)
	if err == nil {
		t.Error("Should not have succeeded")
	}
	if len(missing) != len(LOG2DB_TABLES) {
		t.Error("Should have found all tables missing.")
	}

	err = user.CreateTable(db)
	if err != nil {
		t.Error(err.Error())
	}
	missing, err = CheckTables(db)
	if err == nil {
		t.Error("Should not have succeeded")
	}
	if len(missing) != (len(LOG2DB_TABLES) - 1) {
		t.Error("Should have found all but one tables missing.")
	}
}

func TestCheckTablesComplete(t *testing.T) {
	os.Remove("/tmp/foo.db")

	db, err := global.OpenDatabase("sqlite:///tmp/foo.db")
	defer db.Close()

	if err != nil {
		t.Error(err.Error())
	}

	missing, err := CheckTables(db)
	if err == nil {
		t.Error("Should not have succeeded")
	}
	if len(missing) != len(LOG2DB_TABLES) {
		t.Error("Should have found all tables missing.")
	}

	err = agent.CreateTable(db)
	if err != nil {
		t.Error(err.Error())
	}

	err = destination.CreateTable(db)
	if err != nil {
		t.Error(err.Error())
	}

	err = source.CreateTable(db)
	if err != nil {
		t.Error(err.Error())
	}

	err = user.CreateTable(db)
	if err != nil {
		t.Error(err.Error())
	}

	err = transaction.CreateTable(db)
	if err != nil {
		t.Error(err.Error())
	}

	missing, err = CheckTables(db)
	if err != nil {
		t.Error("Should have succeeded")
	}
	if len(missing) != 0 {
		t.Error("Should have found no tables missing.")
	}
}

func TestCreateAllTables(t *testing.T) {
	os.Remove("/tmp/foo.db")

	db, err := global.OpenDatabase("sqlite:///tmp/foo.db")
	defer db.Close()

	if err != nil {
		t.Error(err.Error())
	}

	missing, err := CheckTables(db)
	if err == nil {
		t.Error("Should not have succeeded")
	}
	if len(missing) != len(LOG2DB_TABLES) {
		t.Error("Should have found all tables missing.")
	}

	err = CreateAllTables(db)

	missing, err = CheckTables(db)
	if err != nil {
		t.Error("Should have succeeded")
	}
	if len(missing) != 0 {
		t.Error("Should have found no tables missing.")
	}
}
