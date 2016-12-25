package user

import (
	"errors"
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

func ReadName(d *sqlx.DB, name string) (*User, error) {
	u := &User{}
	var err error
	row := d.QueryRow("SELECT * FROM users WHERE name=?", name)
	if row == nil {
		err = errors.New("User not found")
	} else {
		err = row.Scan(&u.ID, &u.Name)
	}
	return u, err
}

func Read(d *sqlx.DB, id int) (*User, error) {
	u := &User{}
	var err error
	row := d.QueryRow("SELECT * FROM users WHERE id=?", id)
	if row == nil {
		err = errors.New("User not found")
	} else {
		err = row.Scan(&u.ID, &u.Name)
	}
	return u, err
}

func ReadAll(d *sqlx.DB) ([]*User, error) {
	us := make([]*User, 0)
	rows, err := d.Query("SELECT * FROM users")
	if err == nil {
		for rows.Next() {
			u := &User{}
			rows.Scan(&u.ID, &u.Name)
			if u.ID >= 0 && len(u.Name) > 0 {
				us = append(us, u)
			}
		}
		rows.Close()
	}
	return us, err
}

func Update(d *sqlx.DB, u *User) error {
	_, err := d.Exec("UPDATE users SET name=? WHERE id=?", u.Name, u.ID)
	return err
}
