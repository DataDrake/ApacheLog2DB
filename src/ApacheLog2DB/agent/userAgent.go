package agent

import (
	"database/sql"
	"errors"
)

type UserAgent struct {
	ID   int
	Name string
}

func NewAgent(name string) *UserAgent {
	return &UserAgent{-1, name}
}

func ReadOrCreate(db *sql.DB, name string) (*UserAgent, error) {
	agent, err := ReadName(db, name)
	if err != nil {
		err = Insert(db, NewAgent(name))
		if err == nil {
			agent, err = ReadName(db, name)
		}
	}
	return agent, err
}

func CreateTable(d *sql.DB) error {
	_, err := d.Exec("CREATE TABLE user_agents ( id INTEGER PRIMARY KEY AUTOINCREMENT, name TEXT )")
	return err
}

func Insert(d *sql.DB, u *UserAgent) error {
	_, err := d.Exec("INSERT INTO user_agents VALUES( NULL , ? )", u.Name)
	return err
}

func ReadName(d *sql.DB, name string) (*UserAgent, error) {
	var u *UserAgent
	var err error
	row := d.QueryRow("SELECT * FROM user_agents WHERE name=?", name)
	if row == nil {
		err = errors.New("Agent not found")
	} else {
		var id int
		var name string
		err = row.Scan(&id,&name)
		if err == nil {
			u = &UserAgent{id, name}
		}
	}
	return u, err
}

func Read(d *sql.DB, id int) (*UserAgent, error) {
	var u *UserAgent
	var err error
	row := d.QueryRow("SELECT * FROM user_agents WHERE id=?", id)
	if row == nil {
		err = errors.New("Agent not found")
	} else {
		var id int
		var name string
		err = row.Scan(&id, &name)
		if err == nil {
			u = &UserAgent{id, name}
		}
	}
	return u, err
}

func ReadAll(d *sql.DB) ([]*UserAgent, error) {
	u := make([]*UserAgent, 0)
	tx, err := d.Begin()
	if err == nil {
		rows, err := tx.Query("SELECT * FROM user_agents")
		if err != nil {
			tx.Rollback()
		} else {
			var id int
			var name string
			rows.Scan(&id, &name)
			if id >= 0 && len(name) > 0 {
				u = append(u, &UserAgent{id, name})
			}
			for rows.Next() {
				rows.Scan(&id, &name)
				if id >= 0 && len(name) > 0 {
					u = append(u, &UserAgent{id, name})
				}
			}
			rows.Close()
			tx.Commit()
		}
	}
	return u, err
}

func Update(d *sql.DB, u *UserAgent) error {
	_, err := d.Query("UPDATE user_agents SET name=? WHERE id=?", u.Name, u.ID)
	return err
}
