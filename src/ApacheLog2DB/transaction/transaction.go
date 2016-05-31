package transaction

import (
	"ApacheLog2DB/agent"
	"ApacheLog2DB/destination"
	"ApacheLog2DB/source"
	"ApacheLog2DB/user"
	"database/sql"
	"time"
	"errors"
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
	_, err := d.Exec("CREATE TABLE txns" +
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
	return err
}

func Insert(d *sql.DB, t *Transaction) error {
	_, err := d.Exec("INSERT INTO txns VALUES( NULL,?,?,?,?,?,?,?,?,?,?,? )",
		t.Ident, t.Verb, t.Protocol, t.Status, t.Size, t.Referrer,t.Occurred,
		t.Source.ID,
		t.Dest.ID,
		t.Agent.ID,
		t.User.ID)
	return err
}

func Read(d *sql.DB, id int) (*Transaction, error) {
	var t *Transaction
	var err error
	row := d.QueryRow("SELECT * FROM txns WHERE id=?", id)
	if row != nil {
		err = errors.New("Transaction not found")
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
		err = row.Scan(&id, &ident, &verb, &protocol, &status, &size, &referrer, &occurred, &sourceid, &destid, &agentid, &userid)
		if err == nil {
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
	rows, err := db.Query("SELECT * FROM txns")
	if err == nil {
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
	}
	return t, err
}

func Update(d *sql.DB, u *Transaction) error {
	_, err := d.Query("UPDATE txns SET ident=?,verb=?,protocol=?,status=?,size=?,referrer=?,occurred=?,sourceid=?,destid=?,agentid=?,userid=? WHERE id=?",
		u.Ident, u.Verb, u.Protocol, u.Status, u.Size, u.Referrer, u.Occurred, u.Source.ID, u.Dest.ID, u.Agent.ID, u.User.ID)
	return err
}
