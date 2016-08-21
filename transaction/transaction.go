package transaction

import (
	"database/sql"
	"errors"
	"github.com/DataDrake/ApacheLog2DB/agent"
	"github.com/DataDrake/ApacheLog2DB/destination"
	"github.com/DataDrake/ApacheLog2DB/global"
	"github.com/DataDrake/ApacheLog2DB/source"
	"github.com/DataDrake/ApacheLog2DB/user"
	"time"
	"fmt"
	"os"
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

var CREATE_TABLE = map[string]string{
	"mysql": `CREATE TABLE txns ( id INTEGER AUTO_INCREMENT,
	ident TEXT, verb TEXT, protocol TEXT, status INTEGER,
	size INTEGER, referrer TEXT, occured DATETIME, sourceid INTEGER,
	destid INTEGER, agentid INTEGER, userid INTEGER,
	PRIMARY KEY (id), FOREIGN KEY(sourceid) REFERENCES sources(id),
	FOREIGN KEY(destid) REFERENCES destinations(id),
	FOREIGN KEY(agentid) REFERENCES user_agents(id),
	FOREIGN KEY(userid) REFERENCES users(id) )`,

	"sqlite": `CREATE TABLE txns ( id INTEGER PRIMARY KEY AUTOINCREMENT,
	ident TEXT, verb TEXT, protocol TEXT, status INTEGER,
	size INTEGER, referrer TEXT, occured DATETIME, sourceid INTEGER,
	destid INTEGER, agentid INTEGER, userid INTEGER,
	FOREIGN KEY(sourceid) REFERENCES sources(id),
	FOREIGN KEY(destid) REFERENCES destinations(id),
	FOREIGN KEY(agentid) REFERENCES user_agents(id),
	FOREIGN KEY(userid) REFERENCES users(id) )`,
}

func CreateTable(d *sql.DB) error {
	_, err := d.Exec(CREATE_TABLE[global.DB_TYPE])
	return err
}

func Insert(d *sql.DB, t *Transaction) error {
	_, err := d.Exec("INSERT INTO txns VALUES( NULL,?,?,?,?,?,?,?,?,?,?,? )",
		t.Ident, t.Verb, t.Protocol, t.Status, t.Size, t.Referrer, t.Occurred,
		t.Source.ID,
		t.Dest.ID,
		t.Agent.ID,
		t.User.ID)
	return err
}

func Read(d *sql.DB, id int) (*Transaction, error) {
	t := &Transaction{}
	var err error
	row := d.QueryRow("SELECT * FROM txns WHERE id=?", id)
	if row != nil {
		err = errors.New("Transaction not found")
	} else {
		var sourceid int
		var destid int
		var agentid int
		var userid int
		err = row.Scan(&t.ID, &t.Ident, &t.Verb, &t.Protocol, &t.Status, &t.Size, &t.Referrer, &t.Occurred, &sourceid, &destid, &agentid, &userid)
		if err == nil {
			t.Source, err = source.Read(d, sourceid)
			if err != nil {
				return nil, err
			}

			t.Dest, err = destination.Read(d, destid)
			if err != nil {
				return nil, err
			}

			t.Agent, err = agent.Read(d, agentid)
			if err != nil {
				return nil, err
			}

			t.User, err = user.Read(d, userid)
			if err != nil {
				return nil, err
			}
		}
	}
	return t, err
}

func ReadAll(db *sql.DB) ([]*Transaction, error) {
	ts := make([]*Transaction, 0)
	rows, err := db.Query("SELECT * FROM txns")
	if err == nil {
		for rows.Next() {
			var sourceid int
			var destid int
			var agentid int
			var userid int
			t := &Transaction{}
			rows.Scan(&t.ID, &t.Ident, &t.Verb, &t.Protocol, &t.Status, &t.Size, &t.Referrer, &t.Occurred, &sourceid, &destid, &agentid, &userid)

			t.Source, err = source.Read(db, sourceid)
			if err != nil {
				return nil, err
			}

			t.Dest, err = destination.Read(db, destid)
			if err != nil {
				return nil, err
			}

			t.Agent, err = agent.Read(db, agentid)
			if err != nil {
				return nil, err
			}

			t.User, err = user.Read(db, userid)
			if err != nil {
				return nil, err
			}

			ts = append(ts, t)
		}
		rows.Close()
	}
	return ts, err
}

func Update(d *sql.DB, t *Transaction) error {
	_, err := d.Exec("UPDATE txns SET ident=?,verb=?,protocol=?,status=?,size=?,referrer=?,occurred=?,sourceid=?,destid=?,agentid=?,userid=? WHERE id=?",
		t.Ident, t.Verb, t.Protocol, t.Status, t.Size, t.Referrer, t.Occurred, t.Source.ID, t.Dest.ID, t.Agent.ID, t.User.ID)
	return err
}

func ReadWork(d *sql.DB, sourceid int, start time.Time, stop time.Time) ([]*Transaction, error) {
	ts := make([]*Transaction, 0)
	var err error
	row, err := d.Query("SELECT * FROM txns WHERE sourceid=? AND occured>=? AND occured<?", sourceid, start, stop)
	if err != nil {
		return ts, err
	}
	for row.Next() {
		var sourceid int
		var destid int
		var agentid int
		var userid int
		t := &Transaction{}
		err = row.Scan(&t.ID, &t.Ident, &t.Verb, &t.Protocol, &t.Status, &t.Size, &t.Referrer, &t.Occurred, &sourceid, &destid, &agentid, &userid)
		if err == nil {
			t.Source, err = source.Read(d, sourceid)
			if err != nil {
				fmt.Fprint(os.Stderr,"Could not get source")
				continue
			}

			t.Dest, err = destination.Read(d, destid)
			if err != nil {
				fmt.Fprint(os.Stderr,"Could not get dest")
				continue
			}

			t.Agent, err = agent.Read(d, agentid)
			if err != nil {
				fmt.Fprint(os.Stderr,"Could not get agent")
				continue
			}

			t.User, err = user.Read(d, userid)
			if err != nil {
				fmt.Fprint(os.Stderr,"Could not get user")
				continue
			}
			ts = append(ts, t)
		} else {
			fmt.Println(err.Error())
		}
	}
	fmt.Println(ts)
	return ts, nil
}
