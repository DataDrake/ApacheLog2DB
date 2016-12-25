package agent

import (
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
	_, err := d.Exec("INSERT INTO user_agents VALUES( NULL , $1 )", u.Name)
	return err
}

func ReadName(d *sqlx.DB, name string) (u *UserAgent, err error) {
	u = &UserAgent{}
	err = d.Get(u, "SELECT * FROM user_agents WHERE name=$1", name)
	return
}

func Read(d *sqlx.DB, id int) (u *UserAgent, err error) {
	u = &UserAgent{}
	err = d.Get(u, "SELECT * FROM user_agents WHERE id=$1", id)
	return
}

func ReadAll(d *sqlx.DB) (us []*UserAgent, err error) {
	us = []*UserAgent{}
	err = d.Select(&us, "SELECT * FROM user_agents")
	return us, err
}

func Update(d *sqlx.DB, u *UserAgent) error {
	_, err := d.Exec("UPDATE user_agents SET name=$1 WHERE id=$2", u.Name, u.ID)
	return err
}
