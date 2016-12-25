package agent

import (
	"errors"
	"fmt"
	"github.com/DataDrake/ApacheLog2DB/global"
    "github.com/jmoiron/sqlx"
	"os"
)

type UserAgent struct {
	ID   int
	Name string
}

func NewAgent(name string) *UserAgent {
	return &UserAgent{-1, name}
}

func ReadOrCreate(db *sqlx.DB, name string) (*UserAgent, error) {
	agent, err := ReadName(db, name)
	if err != nil {
		err = Insert(db, NewAgent(name))
		if err != nil {
			fmt.Fprintf(os.Stderr, "[AGENT] Error: %s\n", err.Error())
		}
		agent, err = ReadName(db, name)
	}
	return agent, err
}

var CREATE_TABLE = map[string]string{
	"mysql":  "CREATE TABLE user_agents ( id INTEGER NOT NULL AUTO_INCREMENT, name TEXT , PRIMARY KEY (id) )",
	"sqlite": "CREATE TABLE user_agents ( id INTEGER PRIMARY KEY AUTOINCREMENT, name TEXT )",
}

func CreateTable(d *sqlx.DB) error {
	_, err := d.Exec(CREATE_TABLE[global.DB_TYPE])
	return err
}

func Insert(d *sqlx.DB, u *UserAgent) error {
	_, err := d.Exec("INSERT INTO user_agents VALUES( NULL , ? )", u.Name)
	return err
}

func ReadName(d *sqlx.DB, name string) (*UserAgent, error) {
	u := &UserAgent{}
	var err error
	row := d.QueryRow("SELECT * FROM user_agents WHERE name=?", name)
	if row == nil {
		err = errors.New("Agent not found")
	} else {
		err = row.Scan(&u.ID, &u.Name)
	}
	return u, err
}

func Read(d *sqlx.DB, id int) (*UserAgent, error) {
	u := &UserAgent{}
	var err error
	row := d.QueryRow("SELECT * FROM user_agents WHERE id=?", id)
	if row == nil {
		err = errors.New("Agent not found")
	} else {
		err = row.Scan(&u.ID, &u.Name)
	}
	return u, err
}

func ReadAll(d *sqlx.DB) ([]*UserAgent, error) {
	us := make([]*UserAgent, 0)
	rows, err := d.Query("SELECT * FROM user_agents")
	if err == nil {
		for rows.Next() {
			u := &UserAgent{}
			rows.Scan(&u.ID, &u.Name)
			if u.ID >= 0 && len(u.Name) > 0 {
				us = append(us, u)
			}
		}
		rows.Close()
	}
	return us, err
}

func Update(d *sqlx.DB, u *UserAgent) error {
	_, err := d.Exec("UPDATE user_agents SET name=? WHERE id=?", u.Name, u.ID)
	return err
}
