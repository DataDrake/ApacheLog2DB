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

package source

import (
	"github.com/DataDrake/ApacheLog2DB/global"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
	"testing"
)

func TestTableCreate(t *testing.T) {
	os.Remove("/tmp/foo.db")
	db, err := global.OpenDatabase("sqlite:///tmp/foo.db")
	if err != nil {
		log.Fatal(err)
	}
	err = CreateTable(db)
	if err != nil {
		t.Error(err.Error())
	}
	defer db.Close()
}

func TestReadAllEmpty(t *testing.T) {
	os.Remove("/tmp/foo.db")
	db, err := global.OpenDatabase("sqlite:///tmp/foo.db")
	if err != nil {
		log.Fatal(err)
	}
	err = CreateTable(db)
	if err != nil {
		t.Error(err.Error())
	}
	sources, err := ReadAll(db)
	if err != nil {
		t.Error(err.Error())
	}
	if len(sources) > 0 {
		t.Error("Table should be empty")
	}
	defer db.Close()
}

func TestInsertOne(t *testing.T) {
	os.Remove("/tmp/foo.db")
	db, err := global.OpenDatabase("sqlite:///tmp/foo.db")
	if err != nil {
		log.Fatal(err)
	}
	err = CreateTable(db)
	if err != nil {
		t.Error(err.Error())
	}
	sources, err := ReadAll(db)
	if err != nil {
		t.Error(err.Error())
	}
	if len(sources) > 0 {
		t.Error("Table should be empty")
	}

	err = Insert(db, &Source{0, "192.168.1.1"})
	if err != nil {
		t.Error(err.Error())
	}

	sources, err = ReadAll(db)
	if err != nil {
		t.Error(err.Error())
	}
	if len(sources) != 1 {
		t.Errorf("Table should have 1 entry, found: %d", len(sources))
	}
	defer db.Close()
}
