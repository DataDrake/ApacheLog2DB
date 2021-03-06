//
// Copyright 2016-2017 Bryan T. Meyers <bmeyers@datadrake.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package transaction

import (
	"fmt"
	"github.com/DataDrake/ApacheLog2DB/agent"
	"github.com/DataDrake/ApacheLog2DB/destination"
	"github.com/DataDrake/ApacheLog2DB/global"
	"github.com/DataDrake/ApacheLog2DB/source"
	"github.com/DataDrake/ApacheLog2DB/user"
	"github.com/jmoiron/sqlx"
	"os"
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
	SourceID int
	Dest     *destination.Destination
	DestID   int
	Agent    *agent.UserAgent
	AgentID  int
	User     *user.User
	UserID   int
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

func CreateTable(d *sqlx.DB) error {
	_, err := d.Exec(CREATE_TABLE[global.DB_TYPE])
	return err
}

func Insert(d *sqlx.DB, t *Transaction) error {
	_, err := d.Exec("INSERT INTO txns VALUES( NULL,?,?,?,?,?,?,?,?,?,?,? )",
		t.Ident, t.Verb, t.Protocol, t.Status, t.Size, t.Referrer, t.Occurred,
		t.Source.ID,
		t.Dest.ID,
		t.Agent.ID,
		t.User.ID)
	return err
}

func Read(d *sqlx.DB, id int) (t *Transaction, err error) {
	t = &Transaction{}
	err = d.Get(t, "SELECT * FROM txns WHERE id=$1", id)
	if err != nil {
		return
	}
	t.Source, err = source.Read(d, t.SourceID)
	if err != nil {
		return
	}

	t.Dest, err = destination.Read(d, t.DestID)
	if err != nil {
		return
	}

	t.Agent, err = agent.Read(d, t.AgentID)
	if err != nil {
		return
	}
	t.User, err = user.Read(d, t.UserID)
	return
}

func ReadAll(db *sqlx.DB) ([]*Transaction, error) {
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

func Update(d *sqlx.DB, t *Transaction) error {
	_, err := d.Exec("UPDATE txns SET ident=?,verb=?,protocol=?,status=?,size=?,referrer=?,occurred=?,sourceid=?,destid=?,agentid=?,userid=? WHERE id=?",
		t.Ident, t.Verb, t.Protocol, t.Status, t.Size, t.Referrer, t.Occurred, t.Source.ID, t.Dest.ID, t.Agent.ID, t.User.ID)
	return err
}

func ReadWork(d *sqlx.DB, sourceid int, start time.Time, stop time.Time) ([]*Transaction, error) {
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
				fmt.Fprint(os.Stderr, "Could not get source")
				continue
			}

			t.Dest, err = destination.Read(d, destid)
			if err != nil {
				fmt.Fprint(os.Stderr, "Could not get dest")
				continue
			}

			t.Agent, err = agent.Read(d, agentid)
			if err != nil {
				fmt.Fprint(os.Stderr, "Could not get agent")
				continue
			}

			t.User, err = user.Read(d, userid)
			if err != nil {
				fmt.Fprint(os.Stderr, "Could not get user")
				continue
			}
			ts = append(ts, t)
		} else {
			fmt.Println(err.Error())
		}
	}
	return ts, nil
}
