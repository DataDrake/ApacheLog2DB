package destination

import (
	"database/sql"
	"errors"
	"github.com/DataDrake/ApacheLog2DB/core"
)

type Destination struct {
	ID  int
	URI string
}

func NewDestination(uri string) *Destination {
	return &Destination{-1, uri}
}

func ReadOrCreate(db *sql.DB, uri string) (*Destination, error) {
	dest, err := ReadURI(db, uri)
	if err != nil {
		dest = NewDestination(uri)
		err = Insert(db, dest)
		dest, err = ReadURI(db, uri)
	}
	return dest, err
}

var CREATE_TABLE = map[string]string {
	"mysql":"CREATE TABLE destinations ( id INTEGER AUTO_INCREMENT, uri TEXT, PRIMARY KEY (id))",
	"sqlite":"CREATE TABLE destinations ( id INTEGER AUTOINCREMENT, uri TEXT, PRIMARY KEY (id))",
}

func CreateTable(db *sql.DB) error {
	_, err := db.Exec(CREATE_TABLE[core.DB_TYPE])
	return err
}

func Insert(db *sql.DB, d *Destination) error {
	_, err := db.Exec("INSERT INTO destinations VALUES( NULL , ? )", d.URI)
	return err
}

func ReadURI(db *sql.DB, uri string) (*Destination, error) {
	d := &Destination{}
	var err error
	row := db.QueryRow("SELECT * FROM destinations WHERE uri=?", uri)
	if row == nil {
		err = errors.New("Destination not found")
	} else {
		err = row.Scan(&d.ID, &d.URI)
	}
	return d, err
}

func Read(db *sql.DB, id int) (*Destination, error) {
	d := &Destination{}
	var err error
	row := db.QueryRow("SELECT * FROM destinations WHERE id=?", id)
	if row == nil {
		err = errors.New("Destination not found")
	} else {
		err = row.Scan(&d.ID, &d.URI)
	}
	return d, err
}

func ReadAll(d *sql.DB) ([]*Destination, error) {
	ds := make([]*Destination, 0)
	rows, err := d.Query("SELECT * FROM destinations")
	if err == nil {
		for rows.Next() {
			d := &Destination{}
			rows.Scan(&d.ID, &d.URI)
			if d.ID >= 0 && len(d.URI) > 0 {
				ds = append(ds, d)
			}
		}
		rows.Close()
	}
	return ds, err
}

func Update(db *sql.DB, d *Destination) error {
	_, err := db.Exec("UPDATE destinations SET uri=? WHERE id=?", d.URI, d.ID)
	return err
}
