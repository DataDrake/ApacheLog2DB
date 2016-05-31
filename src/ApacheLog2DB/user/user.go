package user

import "database/sql"

type User struct {
	ID   int
	Name string
}

func CreateTable(d *sql.DB) error {
	tx, err := d.Begin()
	if err == nil {
		_, err = tx.Exec("CREATE TABLE users ( id INTEGER PRIMARY KEY AUTOINCREMENT, name TEXT )")
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}
	return err
}

func Insert(d *sql.DB, u *User) error {
	tx, err := d.Begin()
	if err == nil {
		_, err = tx.Exec("INSERT INTO users VALUES( NULL , ? )", u.Name)
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}
	return err
}

func ReadName(d *sql.DB, name string) (*User, error) {
	var u *User
	tx, err := d.Begin()
	if err == nil {
		row := tx.QueryRow("SELECT * FROM users WHERE name=?", name)
		if row != nil {
			tx.Rollback()
		} else {
			tx.Commit()
			var id int
			var name string
			row.Scan(&id, &name)
			u = &User{id, name}
		}
	}
	return u, err
}

func Read(d *sql.DB, id int) (*User, error) {
	var u *User
	tx, err := d.Begin()
	if err == nil {
		row := tx.QueryRow("SELECT * FROM users WHERE id=?", int)
		if row != nil {
			tx.Rollback()
		} else {
			tx.Commit()
			var id int
			var name string
			row.Scan(&id, &name)
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
	tx, err := d.Begin()
	if err == nil {
		_, err = tx.Query("UPDATE users SET name=? WHERE id=?", u.Name, u.ID)
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}
	return err
}
