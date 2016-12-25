package user

import (
	"github.com/DataDrake/ApacheLog2DB/global"
    "github.com/jmoiron/sqlx"
)

type User struct {
	ID   int
	Name string
}

func NewUser(name string) *User {
	return &User{-1, name}
}

func ReadOrCreate(db *sqlx.DB, name string) (*User, error) {
	src, err := ReadName(db, name)
	if err != nil {
		src = NewUser(name)
		err = Insert(db, src)
		if err == nil {
			src, err = ReadName(db, name)
		}
	}
	return src, err
}

var CREATE_TABLE = map[string]string{
	"mysql":  "CREATE TABLE users ( id INTEGER AUTO_INCREMENT, name TEXT, PRIMARY KEY (id))",
	"sqlite": "CREATE TABLE users ( id INTEGER PRIMARY KEY AUTOINCREMENT, name TEXT)",
}

func CreateTable(d *sqlx.DB) error {
	_, err := d.Exec(CREATE_TABLE[global.DB_TYPE])
	return err
}

func Insert(d *sqlx.DB, u *User) error {
	_, err := d.Exec("INSERT INTO users VALUES( NULL , ? )", u.Name)
	return err
}

func ReadName(d *sqlx.DB, name string) (u *User, err error) {
	u = &User{}
	err = d.Get(u,"SELECT * FROM users WHERE name=$1", name)
	return
}

func Read(d *sqlx.DB, id int) (u *User, err error) {
	u  = &User{}
	err = d.Get(u, "SELECT * FROM users WHERE id=$1", id)
	return
}

func ReadAll(d *sqlx.DB) (us []*User, err error) {
	us = []*User{}
	err = d.Select(&us, "SELECT * FROM users")
	return
}

func Update(d *sqlx.DB, u *User) error {
	_, err := d.Exec("UPDATE users SET name=? WHERE id=?", u.Name, u.ID)
	return err
}
