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

package destination

import (
	"github.com/DataDrake/ApacheLog2DB/global"
	"github.com/jmoiron/sqlx"
)

type Destination struct {
	ID  int
	URI string
}

func NewDestination(uri string) *Destination {
	return &Destination{-1, uri}
}

func ReadOrCreate(db *sqlx.DB, uri string) (*Destination, error) {
	dest, err := ReadURI(db, uri)
	if err != nil {
		dest = NewDestination(uri)
		err = Insert(db, dest)
		dest, err = ReadURI(db, uri)
	}
	return dest, err
}

var CREATE_TABLE = map[string]string{
	"mysql":  "CREATE TABLE destinations ( id INTEGER AUTO_INCREMENT, uri TEXT, PRIMARY KEY (id))",
	"sqlite": "CREATE TABLE destinations ( id INTEGER PRIMARY KEY AUTOINCREMENT, uri TEXT)",
}

func CreateTable(db *sqlx.DB) error {
	_, err := db.Exec(CREATE_TABLE[global.DB_TYPE])
	return err
}

func Insert(db *sqlx.DB, d *Destination) error {
	_, err := db.Exec("INSERT INTO destinations VALUES( NULL , $1 )", d.URI)
	return err
}

func ReadURI(db *sqlx.DB, uri string) (d *Destination, err error) {
	d = &Destination{}
	err = db.Get(d, "SELECT * FROM destinations WHERE uri=$1", uri)
	return
}

func Read(db *sqlx.DB, id int) (d *Destination, err error) {
	d = &Destination{}
	err = db.Get(d, "SELECT * FROM destinations WHERE id=$1", id)
	return
}

func ReadAll(d *sqlx.DB) (ds []*Destination, err error) {
	ds = []*Destination{}
	err = d.Select(&ds, "SELECT * FROM destinations")
	return
}

func Update(db *sqlx.DB, d *Destination) error {
	_, err := db.Exec("UPDATE destinations SET uri=? WHERE id=?", d.URI, d.ID)
	return err
}
