package categories

import (
	"fmt"
	"github.com/joune/zenly/geo"
	"github.com/joune/zenly/util"
	"testing"
)

var user1 = uint64(1)

func byUser2() util.GroupedRows {
	rows := make([]map[string]interface{}, 10)
	rows[0] = map[string]interface{}{"user2": 2, "duration": int64(80), "is_night": false, "location1": int8(geo.Home), "location2": int8(geo.Home)}
	rows[1] = map[string]interface{}{"user2": 2, "duration": int64(80), "is_night": false, "location1": int8(geo.Home), "location2": int8(geo.Home)}

	rows[2] = map[string]interface{}{"user2": 3, "duration": int64(40), "is_night": false, "location1": int8(geo.Other), "location2": int8(geo.Home)}
	rows[3] = map[string]interface{}{"user2": 3, "duration": int64(40), "is_night": false, "location1": int8(geo.Other), "location2": int8(geo.Work)}
	rows[4] = map[string]interface{}{"user2": 3, "duration": int64(40), "is_night": false, "location1": int8(geo.Other), "location2": int8(geo.Other)}

	rows[5] = map[string]interface{}{"user2": 4, "duration": int64(20), "is_night": true, "location1": int8(geo.Home), "location2": int8(geo.Other)}
	rows[6] = map[string]interface{}{"user2": 4, "duration": int64(20), "is_night": true, "location1": int8(geo.Other), "location2": int8(geo.Home)}
	rows[7] = map[string]interface{}{"user2": 4, "duration": int64(20), "is_night": true, "location1": int8(geo.Home), "location2": int8(geo.Other)}

	rows[8] = map[string]interface{}{"user2": 5, "duration": int64(60), "is_night": true, "location1": int8(geo.Home), "location2": int8(geo.Other)}
	rows[9] = map[string]interface{}{"user2": 5, "duration": int64(60), "is_night": true, "location1": int8(geo.Other), "location2": int8(geo.Home)}
	return util.ToRows(rows).GroupByKey("user2")
}

func TestMostSeen(t *testing.T) {
	ms := MostSeen(user1, byUser2())
	if ms != 2 {
		t.Fatal(fmt.Sprintf("MostSeen found %v instead of 2", ms))
	}
}

func TestBestFriend(t *testing.T) {
	bf := BestFriend(user1, byUser2())
	if bf != 3 {
		t.Fatal(fmt.Sprintf("BestFriend found %v instead of 3", bf))
	}
}

func TestCrush(t *testing.T) {
	cr := Crush(user1, byUser2())
	if cr != 4 {
		t.Fatal(fmt.Sprintf("Crush found %v instead of 4", cr))
	}
}
