package categories

import (
	"github.com/gocql/gocql"
	"github.com/joune/zenly/db"
	"github.com/joune/zenly/geo"
	"github.com/joune/zenly/util"
	"strconv"
)

func MostSeen(user1 uint64, byUser2 util.GroupedRows) uint64 {
	user, _ := maxOf(user1, byUser2, func(rows util.Rows) int64 {
		return rows.SumOfKey("duration")
	})
	return user
}

func BestFriend(user1 uint64, byUser2 util.GroupedRows) uint64 {
	return MostSeen(user1, byUser2.FilterGroups(func(row util.Row) bool {
		return row["location1"].(int8) == int8(geo.Other) // outside of user1's home|work
	}))
}

func Crush(user1 uint64, byUser2 util.GroupedRows) uint64 {
	wildNights := byUser2.FilterGroups(func(row util.Row) bool {
		return row["is_night"].(bool) &&
			(row["location1"].(int8) == int8(geo.Home) || //Home but not SameHome!
				row["location2"].(int8) == int8(geo.Home))
	})
	user, nbNights := maxOf(user1, wildNights, func(rows util.Rows) int64 {
		return int64(len(rows))
	})
	if nbNights >= 3 { //at least 3 nights
		return user
	} else {
		return user1 //sorry, not wild enough!
	}
}

func MutualLove_7Days(cql *gocql.Session, user1 uint64, user2 uint64) uint64 {
	if user2 == user1 {
		return user1 //user2 not found, spare the additional work
	}
	byUser2, err := db.FetchSessionsGrouped(cql, user2, 7)
	if err != nil {
		return user1 //not found
	}
	if (MostSeen(user2, byUser2)) == user1 {
		return user2 //MostSeen for both
	} else {
		return user1 //no match
	}
}

func MutualLoveGlobal(cql *gocql.Session, user1 uint64) uint64 {
	user2 := mostSeenGlobal(cql, user1)
	if user2 == user1 {
		return user1 //user2 not found, spare the additional work
	}
	if (mostSeenGlobal(cql, user2)) == user1 {
		return user2 //MostSeen for both
	} else {
		return user1 //no match
	}
}

func mostSeenGlobal(cql *gocql.Session, user1 uint64) uint64 {
	durations, err := db.FetchDurations(cql, user1)
	if err != nil {
		return user1 //not found
	}
	user2, duration := user1, int64(0)
	for i := 0; i < len(durations); i++ {
		if d := durations[i]["duration"].(int64); d > duration {
			user2, duration = uint64(durations[i]["user2"].(int64)), d //wtf?
		}
	}
	return user2
}

// scan rows grouped by user and yield (userId,maxCount) when applying the counter function
// to a group of rows. Default to user1 if no other maximum is found
func maxOf(user1 uint64, byUser2 util.GroupedRows, counter func(util.Rows) int64) (uint64, int64) {
	// see GroupByKey about awkward conversion to string
	// don't fire me for this, we can talk it out during the code review :)
	user, max := strconv.FormatUint(user1, 10), int64(0)
	for user2, rows := range byUser2 {
		if count := counter(rows); count > max {
			user, max = user2, count
		}
	}
	res, err := strconv.ParseUint(user, 10, 64)
	if err != nil {
		panic(err)
	}
	return res, max
}
