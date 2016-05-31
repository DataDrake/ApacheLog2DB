package destination

import (
	"database/sql"
	"errors"
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

func CreateTable(db *sql.DB) error {
	_, err:= db.Exec("CREATE TABLE destinations ( id INTEGER PRIMARY KEY AUTOINCREMENT, uri TEXT )")
	return err
}

func Insert(db *sql.DB, d *Destination) error {
	_, err := db.Exec("INSERT INTO destinations VALUES( NULL , ? )", d.URI)
	return err
}

func ReadURI(db *sql.DB, uri string) (*Destination, error) {
	var d *Destination
	var err error
	row := db.QueryRow("SELECT * FROM destinations WHERE uri=?", uri)
	if row == nil {
		err = errors.New("Destination not found")
	} else {
		var id int
		var uri string
		err = row.Scan(&id, &uri)
		if err == nil {
			d = &Destination{id, uri}
		}
	}
	return d, err
}

func Read(db *sql.DB, id int) (*Destination, error) {
	var d *Destination
	var err error
	row := db.QueryRow("SELECT * FROM destinations WHERE id=?", id)
	if row == nil {
		err = errors.New("Destination not found")
	} else {
		var id int
		var uri string
		err = row.Scan(&id, &uri)
		if err == nil {
			d = &Destination{id, uri}
		}
	}
	return d, err
}

func ReadAll(d *sql.DB) ([]*Destination, error) {
	ds := make([]*Destination, 0)
	tx, err := d.Begin()
	if err == nil {
		rows, err := tx.Query("SELECT * FROM destinations")
		if err != nil {
			tx.Rollback()
		} else {
			var id int
			var uri string
			rows.Scan(&id, &uri)
			if id >= 0 && len(uri) > 0 {
				ds = append(ds, &Destination{id, uri})
			}
			for rows.Next() {
				rows.Scan(&id, &uri)
				if id >= 0 && len(uri) > 0 {
					ds = append(ds, &Destination{id, uri})
				}
			}
			rows.Close()
			tx.Commit()
		}
	}
	return ds, err
}

func Update(db *sql.DB, d *Destination) error {
	_, err := db.Query("UPDATE destinations SET uri=? WHERE id=?", d.URI, d.ID)
	return err
}
