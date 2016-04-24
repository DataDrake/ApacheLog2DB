package transaction

import (
	"ApacheLog2DB/agent"
	"ApacheLog2DB/destination"
	"ApacheLog2DB/source"
	"ApacheLog2DB/user"
	"database/sql"
	"time"
)

type Transaction struct {
	ID       int
	Ident    string
	Verb     string
	Protocol string
	Status   int
	Size     int
	Referrer string
	Occurred time.Time
	Source   *source.Source
	Dest     *destination.Destination
	Agent    *agent.UserAgent
	User     *user.User
}

func CreateTable(d *sql.DB) error {
	tx, err := d.Begin()
	if err == nil {
		_, err = tx.Exec("CREATE TABLE transactions" +
			"( id INTEGER PRIMARY KEY AUTOINCREMENT," +
			"ident TEXT, " +
			"verb TEXT, " +
			"protocol TEST, " +
			"status INTEGER, " +
			"size INTEGER, " +
			"referrer TEXT, " +
			"occured DATETIME, " +
			"sourceid INTEGER, " +
			"destid INTEGER, " +
			"agentid INTEGER, " +
			"userid INTEGER, " +
			"FOREIGN KEY(sourceid) REFERENCES sources(id), " +
			"FOREIGN KEY(destid) REFERENCES destinations(id), " +
			"FOREIGN KEY(agentid) REFERENCES user_agents(id), " +
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

func Insert(d *sql.DB, t *Transaction) error {
	tx, err := d.Begin()
	if err == nil {
		_, err = tx.Exec("INSERT INTO transactions VALUES( NULL,?,?,?,?,?,?,?,?,?,?,? )",
			t.Ident, t.Verb, t.Protocol, t.Status, t.Size, t.Referrer,
			t.Source.ID,
			t.Dest.ID,
			t.Agent.ID,
			t.User.ID)
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}
	return err
}

func Read(d *sql.DB, id int) (*Transaction, error) {
	var t *Transaction
	tx, err := d.Begin()
	if err == nil {
		row := tx.QueryRow("SELECT * FROM transaction WHERE id=?", id)
		if row != nil {
			tx.Rollback()
		} else {
			tx.Commit()
			var id int
			var ident string
			var verb string
			var protocol string
			var status int
			var size int
			var referrer string
			var occurred time.Time
			var sourceid int
			var destid int
			var agentid int
			var userid int
			row.Scan(&id, &ident, &verb, &protocol, &status, &size, &referrer, &occurred, &sourceid, &destid, &agentid, &userid)
			source, err := source.Read(d, sourceid)
			if err != nil {
				return nil, err
			}
			dest, err := destination.Read(d, destid)
			if err != nil {
				return nil, err
			}
			agent, err := agent.Read(d, agentid)
			if err != nil {
				return nil, err
			}
			user, err := user.Read(d, userid)
			if err != nil {
				return nil, err
			}
			t = &Transaction{id, ident, verb, protocol, status, size, referrer, occurred, source, dest, agent, user}
		}
	}
	return t, err
}

func ReadAll(db *sql.DB) ([]*Transaction, error) {
	t := make([]*Transaction, 0)
	tx, err := db.Begin()
	if err == nil {
		rows, err := tx.Query("SELECT * FROM transactions")
		if err != nil {
			tx.Rollback()
		} else {
			var id int
			var ident string
			var verb string
			var protocol string
			var status int
			var size int
			var referrer string
			var occurred time.Time
			var sourceid int
			var destid int
			var agentid int
			var userid int
			rows.Scan(&id, &ident, &verb, &protocol, &status, &size, &referrer, &occurred, &sourceid, &destid, &agentid, &userid)
			s, err := source.Read(db, sourceid)
			if err != nil {
				return nil, err
			}
			d, err := destination.Read(db, destid)
			if err != nil {
				return nil, err
			}
			a, err := agent.Read(db, agentid)
			if err != nil {
				return nil, err
			}
			u, err := user.Read(db, userid)
			if err != nil {
				return nil, err
			}
			if id >= 0 {
				t = append(t, &Transaction{id, ident, verb, protocol, status, size, referrer, occurred, s, d, a, u})
			}
			for rows.Next() {
				rows.Scan(&id, &ident, &verb, &protocol, &status, &size, &referrer, &occurred, &sourceid, &destid, &agentid, &userid)
				s, err := source.Read(db, sourceid)
				if err != nil {
					return nil, err
				}
				d, err := destination.Read(db, destid)
				if err != nil {
					return nil, err
				}
				a, err := agent.Read(db, agentid)
				if err != nil {
					return nil, err
				}
				u, err := user.Read(db, userid)
				if err != nil {
					return nil, err
				}

				if id >= 0 {
					t = append(t, &Transaction{id, ident, verb, protocol, status, size, referrer, occurred, s, d, a, u})
				}
			}
			rows.Close()
			tx.Commit()
		}
	}
	return t, err
}

func Update(d *sql.DB, u *Transaction) error {
	tx, err := d.Begin()
	if err == nil {
		_, err = tx.Query("UPDATE transactions SET ident=?,verb=?,protocol=?,status=?,size=?,referrer=?,occurred=?,sourceid=?,destid=?,agentid=?,userid=? WHERE id=?",
			u.Ident, u.Verb, u.Protocol, u.Status, u.Size, u.Referrer, u.Occurred, u.Source.ID, u.Dest.ID, u.Agent.ID, u.User.ID)
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}
	return err
}
