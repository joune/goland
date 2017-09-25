package util

import (
	"fmt"
	"testing"
	"time"
)

func TestIsNight(t *testing.T) {
	//day time
	t1, _ := time.Parse(time.RFC3339, "2017-09-18T09:05:00+00:00")
	t2, _ := time.Parse(time.RFC3339, "2017-09-18T21:05:00+00:00")
	if IsNight(t1, t2) {
		t.Fatal(fmt.Sprintf("%v - %v", t1, t2))
	}
	//night time
	t3, _ := time.Parse(time.RFC3339, "2017-09-18T22:05:00+00:00")
	t4, _ := time.Parse(time.RFC3339, "2017-09-19T07:05:00+00:00")
	if !IsNight(t3, t4) {
		t.Fatal(fmt.Sprintf("%v - %v", t3, t4))
	}
	//more than 24hrs
	t5, _ := time.Parse(time.RFC3339, "2017-09-16T10:05:00+00:00")
	t6, _ := time.Parse(time.RFC3339, "2017-09-20T10:05:00+00:00")
	if !IsNight(t5, t6) {
		t.Fatal(fmt.Sprintf("%v - %v", t5, t6))
	}
	//less than 6hrs at night
	t7, _ := time.Parse(time.RFC3339, "2017-09-18T10:05:00+00:00")
	t8, _ := time.Parse(time.RFC3339, "2017-09-19T02:05:00+00:00")
	if IsNight(t7, t8) {
		t.Fatal(fmt.Sprintf("%v - %v", t7, t8))
	}
	//different timezone
	t9, _ := time.Parse(time.RFC3339, "2017-09-18T22:05:00+12:00")
	t10, _ := time.Parse(time.RFC3339, "2017-09-19T07:05:00+12:00")
	if !IsNight(t9, t10) {
		t.Fatal(fmt.Sprintf("%v - %v", t9, t10))
	}
}

func TestGroupByKey(t *testing.T) {
	rows := make([]map[string]interface{}, 5)
	rows[0] = map[string]interface{}{"a": 1, "b": 2}
	rows[1] = map[string]interface{}{"a": 2, "b": 4}
	rows[2] = map[string]interface{}{"a": 3, "b": 6}
	rows[3] = map[string]interface{}{"a": 2, "b": 8}
	rows[4] = map[string]interface{}{"a": 1, "b": 10}
	groups := ToRows(rows).GroupByKey("a")
	fmt.Printf("%v", groups)
	if len(groups["1"]) != 2 || len(groups["2"]) != 2 || len(groups["3"]) != 1 {
		t.Fatal("grouping failed")
	}
}
