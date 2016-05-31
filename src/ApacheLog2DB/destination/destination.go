package destination

import "database/sql"

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

func CreateTable(d *sql.DB) error {
	tx, err := d.Begin()
	if err == nil {
		_, err = tx.Exec("CREATE TABLE destinations ( id INTEGER PRIMARY KEY AUTOINCREMENT, uri TEXT )")
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}
	return err
}

func Insert(db *sql.DB, d *Destination) error {
	tx, err := db.Begin()
	if err == nil {
		_, err = tx.Exec("INSERT INTO destinations VALUES( NULL , ? )", d.URI)
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}
	return err
}

func ReadURI(db *sql.DB, uri string) (*Destination, error) {
	var d *Destination
	tx, err := db.Begin()
	if err == nil {
		row := tx.QueryRow("SELECT * FROM destinations WHERE uri=?", uri)
		if row != nil {
			tx.Rollback()
		} else {
			tx.Commit()
			var id int
			var uri string
			row.Scan(&id, &uri)
			d = &Destination{id, uri}
		}
	}
	return d, err
}

func Read(db *sql.DB, id int) (*Destination, error) {
	var d *Destination
	tx, err := db.Begin()
	if err == nil {
		row := tx.QueryRow("SELECT * FROM destinations WHERE id=?", id)
		if row != nil {
			tx.Rollback()
		} else {
			tx.Commit()
			var id int
			var uri string
			row.Scan(&id, &uri)
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
	tx, err := db.Begin()
	if err == nil {
		_, err = tx.Query("UPDATE destinations SET uri=? WHERE id=?", d.URI, d.ID)
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}
	return err
}
