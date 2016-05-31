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
	_, err := d.Exec("CREATE TABLE users ( id INTEGER PRIMARY KEY AUTOINCREMENT, name TEXT )")
	return err
}

func Insert(d *sql.DB, u *User) error {
	_, err := d.Exec("INSERT INTO users VALUES( NULL , ? )", u.Name)
	return err
}

func ReadName(d *sql.DB, name string) (*User, error) {
	var u *User
	var err error
	row := d.QueryRow("SELECT * FROM users WHERE name=?", name)
	if row == nil {
		err = errors.New("User not found")
	} else {
		var id int
		var name string
		err = row.Scan(&id, &name)
		if err == nil {
			u = &User{id, name}
		}
	}
	return u, err
}

func Read(d *sql.DB, id int) (*User, error) {
	var u *User
	var err error
	row := d.QueryRow("SELECT * FROM users WHERE id=?", id)
	if row == nil {
		err = errors.New("User not found")
	} else {
		var id int
		var name string
		err = row.Scan(&id, &name)
		if err == nil {
			u = &User{id, name}
		}
	}
	return u, err
}

func ReadAll(d *sql.DB) ([]*User, error) {
	s := make([]*User, 0)
	tx, err := d.Begin()
	if err == nil {
		rows, err := tx.Query("SELECT * FROM users")
		if err != nil {
			tx.Rollback()
		} else {
			var id int
			var ip string
			rows.Scan(&id, &ip)
			if id >= 0 && len(ip) > 0 {
				s = append(s, &User{id, ip})
			}
			for rows.Next() {
				rows.Scan(&id, &ip)
				if id >= 0 && len(ip) > 0 {
					s = append(s, &User{id, ip})
				}
			}
			rows.Close()
			tx.Commit()
		}
	}
	return s, err
}

func Update(d *sql.DB, u *User) error {
	_, err := d.Query("UPDATE users SET name=? WHERE id=?", u.Name, u.ID)
	return err
}
