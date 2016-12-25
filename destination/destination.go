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
