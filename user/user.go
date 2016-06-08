package user

import (
	"database/sql"
	"errors"
)

type User struct {
	ID   int
	Name string
}

func NewUser(name string) *User {
	return &User{-1, name}
}

func ReadOrCreate(db *sql.DB, name string) (*User, error) {
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

func CreateTable(d *sql.DB) error {
	_, err := d.Exec("CREATE TABLE users ( id INTEGER AUTO_INCREMENT, name TEXT, PRIMARY KEY (id))")
	return err
}

func Insert(d *sql.DB, u *User) error {
	_, err := d.Exec("INSERT INTO users VALUES( NULL , ? )", u.Name)
	return err
}

func ReadName(d *sql.DB, name string) (*User, error) {
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

func Read(d *sql.DB, id int) (*User, error) {
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

func ReadAll(d *sql.DB) ([]*User, error) {
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

func Update(d *sql.DB, u *User) error {
	_, err := d.Exec("UPDATE users SET name=? WHERE id=?", u.Name, u.ID)
	return err
}
