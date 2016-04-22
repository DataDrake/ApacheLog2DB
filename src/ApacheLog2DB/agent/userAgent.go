package agent
import "database/sql"

type UserAgent struct{
	id int
	name string
}

func CreateTable(d *sql.DB) error {
	tx,err := d.Begin()
	if err == nil {
		_,err = tx.Exec("CREATE TABLE user_agents ( id INTEGER , name TEXT )")
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}
	return err
}

func Insert(d *sql.DB, u *UserAgent) error {
	tx,err := d.Begin()
	if err == nil {
		_,err = tx.Query("INSERT INTO user_agents VALUES( NULL , ? )",u.name)
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
			agent,err := Read(d,u.name)
			if err == nil {
				u = agent
			}
		}
	}
	return err
}

func Read(d *sql.DB, name string) (*UserAgent,error) {
	var u *UserAgent
	tx,err := d.Begin()
	if err == nil {
		row,err := tx.QueryRow("SELECT * FROM user_agents WHERE id=?",name)
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
			var id int
			var name string
			row.Scan(&id,&name)
			u = &UserAgent{id,name}
		}
	}
	d.Close()
	return u,err
}

func Update(d *sql.DB,u *UserAgent) error {
	tx,err := d.Begin()
	if err == nil {
		_,err = tx.Query("UPDATE user_agents SET name=? WHERE id=?",u.name,u.id)
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}
	d.Close()
	return err
}

