package agent
import "database/sql"

type UserAgent struct{
	ID   int
	Name string
}

func CreateTable(d *sql.DB) error {
	tx,err := d.Begin()
	if err == nil {
		_,err = tx.Exec("CREATE TABLE user_agents ( id INTEGER PRIMARY KEY AUTOINCREMENT, name TEXT )")
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
		_,err = tx.Exec("INSERT INTO user_agents VALUES( NULL , ? )",u.Name)
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}
	return err
}

func Read(d *sql.DB, id int) (*UserAgent,error) {
	var u *UserAgent
	tx,err := d.Begin()
	if err == nil {
		row := tx.QueryRow("SELECT * FROM user_agents WHERE id=?", id)
		if row != nil {
			tx.Rollback()
		} else {
			tx.Commit()
			var id int
			var name string
			row.Scan(&id,&name)
			u = &UserAgent{id,name}
		}
	}
	return u,err
}

func ReadAll(d *sql.DB) ([]*UserAgent,error) {
	u := make([]*UserAgent,0)
	tx,err := d.Begin()
	if err == nil {
		rows,err := tx.Query("SELECT * FROM user_agents")
		if err != nil {
			tx.Rollback()
		} else {
			var id int
			var name string
			rows.Scan(&id,&name)
			if id >=0 && len(name) > 0 {
				u = append(u,&UserAgent{id,name})
			}
			for rows.Next() {
				rows.Scan(&id,&name)
				if id >= 0 && len(name) > 0 {
					u = append(u,&UserAgent{id,name})
				}
			}
			rows.Close()
			tx.Commit()
		}
	}
	return u,err
}

func Update(d *sql.DB,u *UserAgent) error {
	tx,err := d.Begin()
	if err == nil {
		_,err = tx.Query("UPDATE user_agents SET name=? WHERE id=?",u.Name,u.ID)
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}
	return err
}
