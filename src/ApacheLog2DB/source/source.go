package source

import "database/sql"

type Source struct {
	ID int
	IP string
}

func CreateTable(d *sql.DB) error {
	tx, err := d.Begin()
	if err == nil {
		_, err = tx.Exec("CREATE TABLE sources ( id INTEGER PRIMARY KEY AUTOINCREMENT, ip TEXT )")
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}
	return err
}

func Insert(d *sql.DB, s *Source) error {
	tx, err := d.Begin()
	if err == nil {
		_, err = tx.Exec("INSERT INTO sources VALUES( NULL , ? )", s.IP)
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}
	return err
}

func ReadIP(d *sql.DB, ip string) (*Source, error) {
	var s *Source
	tx, err := d.Begin()
	if err == nil {
		row := tx.QueryRow("SELECT * FROM sources WHERE ip=?", ip)
		if row != nil {
			tx.Rollback()
		} else {
			tx.Commit()
			var id int
			var ip string
			row.Scan(&id, &ip)
			s = &Source{id, ip}
		}
	}
	return s, err
}

func Read(d *sql.DB, id int) (*Source, error) {
	var s *Source
	tx, err := d.Begin()
	if err == nil {
		row := tx.QueryRow("SELECT * FROM sources WHERE id=?", id)
		if row != nil {
			tx.Rollback()
		} else {
			tx.Commit()
			var id int
			var ip string
			row.Scan(&id, &ip)
			s = &Source{id, ip}
		}
	}
	return s, err
}

func ReadAll(d *sql.DB) ([]*Source, error) {
	s := make([]*Source, 0)
	tx, err := d.Begin()
	if err == nil {
		rows, err := tx.Query("SELECT * FROM sources")
		if err != nil {
			tx.Rollback()
		} else {
			var id int
			var ip string
			rows.Scan(&id, &ip)
			if id >= 0 && len(ip) > 0 {
				s = append(s, &Source{id, ip})
			}
			for rows.Next() {
				rows.Scan(&id, &ip)
				if id >= 0 && len(ip) > 0 {
					s = append(s, &Source{id, ip})
				}
			}
			rows.Close()
			tx.Commit()
		}
	}
	return s, err
}

func Update(d *sql.DB, s *Source) error {
	tx, err := d.Begin()
	if err == nil {
		_, err = tx.Query("UPDATE sources SET ip=? WHERE id=?", s.IP, s.ID)
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}
	return err
}
