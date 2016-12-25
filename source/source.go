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
