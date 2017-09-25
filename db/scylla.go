package db

import (
	"github.com/gocql/gocql"
	"github.com/joune/zenly/geo"
	"github.com/joune/zenly/util"
	"log"
	"time"
)

const (
	// mono device hypothesis: the user shouldn't be in 2 places at the same time
	createSessionsTable = `CREATE TABLE IF NOT EXISTS zenly.sessions (
		user1 bigint, 
		user2 bigint,
		latitude double,
		longitude double,
		location1 tinyint, 
		location2 tinyint,
		start_date timestamp,
		end_date timestamp,
		duration bigint,
		is_night boolean,
		PRIMARY KEY (user1, start_date))
		WITH CLUSTERING ORDER BY (start_date DESC);`

	createAggregateTable = `CREATE TABLE IF NOT EXISTS zenly.mutual_love (
		user1 bigint,
		user2 bigint,
		duration counter,
		PRIMARY KEY (user1, user2))`

	insertSession = `INSERT INTO sessions 
		(user1, user2, latitude, longitude, location1, location2, start_date, end_date, duration, is_night) 
		values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	selectSessions = `SELECT user2, duration, location1, location2, is_night 
		 FROM sessions 
		 WHERE user1 = ? and start_date > ?`

	updateDuration = `UPDATE mutual_love set duration = duration + ? where user1 = ? and user2 = ?`

	selectDurations = `SELECT user2, duration FROM mutual_love WHERE user1 = ?`
)

func InitDB() *gocql.Session {
	cluster := gocql.NewCluster("scylladb")
	cluster.Keyspace = "zenly"
	var session *gocql.Session = nil
	var err error = nil
	for session, err = cluster.CreateSession(); err != nil; {
		log.Printf("Failed to connect to ScyllaDB. Will retry in 5 seconds. %s", err)
		time.Sleep(5 * time.Second)
	}
	if err := session.Query(createSessionsTable).Exec(); err != nil {
		log.Printf("ERROR! %v\n", err) //FIXME? don't panic! timeout but the table seems to be created anyway!
	}
	if err := session.Query(createAggregateTable).Exec(); err != nil {
		log.Printf("ERROR! %v\n", err) //FIXME? don't panic! timeout but the table seems to be created anyway!
	}
	return session
}

func InsertSession(cql *gocql.Session,
	user1, user2 uint64,
	latitude, longitude float64,
	location1, location2 geo.Location,
	start_date, end_date time.Time,
	duration time.Duration, is_night bool) {
	if err := cql.Query(insertSession,
		user1, user2,
		latitude, longitude,
		int8(location1), int8(location2),
		start_date, end_date,
		duration, is_night).Exec(); err != nil {
		log.Printf("ERROR! %s\n", err)
	}
}

func FetchSessionsGrouped(cql *gocql.Session, user1 uint64, daysAgo int) (util.GroupedRows, error) {
	_daysAgo := time.Now().AddDate(0, 0, -daysAgo)
	// since sessions for a single user over 7 days shouldn't be too many,
	// we could use "allow filtering", or materialized views, or indexes..
	// but none of this is supported by Scylla :(
	rows, err := fetchRows(cql.Query(selectSessions, user1, _daysAgo))
	if err != nil {
		return nil, err
	}
	log.Printf("Found %v session(s) for %v since %v\n", len(rows), user1, _daysAgo)
	return rows.GroupByKey("user2"), nil
}

func fetchRows(query *gocql.Query) (util.Rows, error) {
	iter := query.Iter()
	slice, err := iter.SliceMap() //sessions for one user over 7 days should fit in memory
	iter.Close()
	if err != nil {
		log.Printf("ERROR! %s\n", err)
		return nil, err
	}
	return util.ToRows(slice), nil
}

func UpdateDuration(cql *gocql.Session, user1, user2 uint64, duration time.Duration) {
	if err := cql.Query(updateDuration, int(duration.Minutes()), user1, user2).Exec(); err != nil {
		log.Printf("ERROR! %s\n", err)
	}
}

func FetchDurations(cql *gocql.Session, user1 uint64) (util.Rows, error) {
	return fetchRows(cql.Query(selectDurations, user1))
}
