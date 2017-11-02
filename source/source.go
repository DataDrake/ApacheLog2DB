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
	"github.com/jmoiron/sqlx"
)

type Source struct {
	ID int
	IP string
}

func NewSource(ip string) *Source {
	return &Source{-1, ip}
}

func ReadOrCreate(db *sqlx.DB, IP string) (*Source, error) {
	src, err := ReadIP(db, IP)
	if err != nil {
		err = Insert(db, NewSource(IP))
		if err == nil {
			src, err = ReadIP(db, IP)
		}
	}
	return src, err
}

var CREATE_TABLE = map[string]string{
	"mysql":  "CREATE TABLE sources ( id INTEGER AUTO_INCREMENT, ip TEXT , PRIMARY KEY (id))",
	"sqlite": "CREATE TABLE sources ( id INTEGER PRIMARY KEY AUTOINCREMENT, ip TEXT)",
}

func CreateTable(d *sqlx.DB) error {
	_, err := d.Exec(CREATE_TABLE[global.DB_TYPE])
	return err
}

func Insert(d *sqlx.DB, s *Source) error {
	_, err := d.Exec("INSERT INTO sources VALUES( NULL , $1 )", s.IP)
	return err
}

func ReadIP(d *sqlx.DB, ip string) (s *Source, err error) {
	s = &Source{}
	err = d.Get(s, "SELECT * FROM sources WHERE ip=$1", ip)
	return
}

func Read(d *sqlx.DB, id int) (s *Source, err error) {
	s = &Source{}
	err = d.Get(s, "SELECT * FROM sources WHERE id=$1", id)
	return
}

func ReadAll(d *sqlx.DB) (ss []*Source, err error) {
	ss = []*Source{}
	err = d.Select(&ss, "SELECT * FROM sources")
	return
}

func Update(d *sqlx.DB, s *Source) error {
	_, err := d.Exec("UPDATE sources SET ip=? WHERE id=?", s.IP, s.ID)
	return err
}
