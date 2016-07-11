package source

import (
	"database/sql"
	"errors"
	"github.com/DataDrake/ApacheLog2DB/core"
)

type Source struct {
	ID int
	IP string
}

func NewSource(ip string) *Source {
	return &Source{-1, ip}
}

func ReadOrCreate(db *sql.DB, IP string) (*Source, error) {
	src, err := ReadIP(db, IP)
	if err != nil {
		err = Insert(db, NewSource(IP))
		if err == nil {
			src, err = ReadIP(db, IP)
		}
	}
	return src, err
}

var CREATE_TABLE = map[string]string {
	"mysql":"CREATE TABLE sources ( id INTEGER AUTO_INCREMENT, ip TEXT , PRIMARY KEY (id))",
	"sqlite":"CREATE TABLE sources ( id INTEGER AUTOINCREMENT, ip TEXT , PRIMARY KEY (id))",
}

func CreateTable(d *sql.DB) error {
	_, err := d.Exec(CREATE_TABLE[core.DB_TYPE])
	return err
}

func Insert(d *sql.DB, s *Source) error {
	_, err := d.Exec("INSERT INTO sources VALUES( NULL , ? )", s.IP)
	return err
}

func ReadIP(d *sql.DB, ip string) (*Source, error) {
	s := &Source{}
	var err error
	row := d.QueryRow("SELECT * FROM sources WHERE ip=?", ip)
	if row == nil {
		err = errors.New("Source not found")
	} else {
		err = row.Scan(&s.ID, &s.IP)
	}
	return s, err
}

func Read(d *sql.DB, id int) (*Source, error) {
	s := &Source{}
	var err error
	row := d.QueryRow("SELECT * FROM sources WHERE id=?", id)
	if row == nil {
		err = errors.New("Source not found")
	} else {
		err = row.Scan(&s.ID, &s.IP)
	}
	return s, err
}

func ReadAll(d *sql.DB) ([]*Source, error) {
	ss := make([]*Source, 0)
	rows, err := d.Query("SELECT * FROM sources")
	if err == nil {
		for rows.Next() {
			s := &Source{}
			rows.Scan(&s.ID, &s.IP)
			if s.ID >= 0 && len(s.IP) > 0 {
				ss = append(ss, s)
			}
		}
		rows.Close()
	}
	return ss, err
}

func Update(d *sql.DB, s *Source) error {
	_, err := d.Exec("UPDATE sources SET ip=? WHERE id=?", s.IP, s.ID)
	return err
}
