package transaction
import (
	"database/sql"
	"time"
	"ApacheLog2DB/source"
	"ApacheLog2DB/agent"
	"ApacheLog2DB/destination"
)

type Transaction struct{
	id int
	ident string
	verb string
	protocol string
	status int
	size int
	referrer string
	occurred time.Time
	source *source.Source
	dest *destination.Destination
	agent *agent.UserAgent
	user *user.User
}

func CreateTable(d *sql.DB) error {
	tx,err := d.Begin()
	if err == nil {
		_,err = tx.Exec("CREATE TABLE transactions" +
						"( id INTEGER PRIMARY KEY AUTOINCREMENT," +
						"ident TEXT, " +
						"verb TEXT, " +
						"protocol TEST, " +
						"status INTEGER, " +
						"size INTEGER, " +
						"referrer TEXT, " +
						"occured DATETIME, " +
						"sourceid INTEGER, " +
						"FOREIGN KEY(sourceid) REFERENCES sources(id), " +
						"destid INTEGER, " +
						"FOREIGN KEY(destid) REFERENCES destinations(id), " +
						"agentid INTEGER, " +
						"FOREIGN KEY(agentid) REFERENCES user_agents(id), " +
						"userid INTEGER, " +
						"FOREIGN KEY(userid) REFERENCES users(id)" +
						" )")
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
		_,err = tx.Exec("INSERT INTO user_agents VALUES( NULL , ? )",u.name)
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}
	return err
}

func Read(d *sql.DB, name string) (*UserAgent,error) {
	var u *UserAgent
	tx,err := d.Begin()
	if err == nil {
		row := tx.QueryRow("SELECT * FROM user_agents WHERE id=?",name)
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
		_,err = tx.Query("UPDATE user_agents SET name=? WHERE id=?",u.name,u.id)
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}
	return err
}
