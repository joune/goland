package util

import (
	"fmt"
	"time"
)

// compute night time with respect to session start date
func IsNight(sStart time.Time, sEnd time.Time) bool {
	nightStart := time.Date(sStart.Year(), sStart.Month(), sStart.Day(), 22, 0, 0, 0, sStart.Location())
	nightEnd := time.Date(nightStart.Year(), nightStart.Month(), nightStart.Day()+1, 8, 0, 0, 0, nightStart.Location())
	nightDuration := Min(nightEnd.Unix(), sEnd.Unix()) - Max(sStart.Unix(), nightStart.Unix())
	return nightDuration >= (3600 * 6)
}

// see https://mrekucci.blogspot.fr/2015/07/dont-abuse-mathmax-mathmin.html
// to "justify" why Go shouldn't have Min/Max for ints
func Min(a, b int64) int64 {
	if a <= b {
		return a
	} else {
		return b
	}
}
func Max(a, b int64) int64 {
	if a >= b {
		return a
	} else {
		return b
	}
}

type Row map[string]interface{}
type Rows []Row
type GroupedRows map[string][]Row

// Is this starting to be absurd ?
// the idea was to make the GroupByKey function work on my Rows type (just to play with types)
// but when reading the slice from DB I get
// "cannot convert slice (type []map[string]interface {}) to type util.Rows"
// so.. FBI? (Fausse Bonne Idée?)
func ToRows(rows []map[string]interface{}) Rows {
	res := make(Rows, len(rows))
	for i := 0; i < len(rows); i++ {
		res[i] = Row(rows[i])
	}
	return res
}

// no generics :/
// so I can either treat everything as a string and parse the result as needed
// or implement it specifically for an expected value type..
// honestly I'm not sure what to choose! (but I did!)
func (rows Rows) GroupByKey(k string) GroupedRows {
	m := make(GroupedRows)
	return _groupByKey(k, rows, m)
}

func _groupByKey(k string, rows []Row, agg GroupedRows) GroupedRows {
	if len(rows) == 0 {
		return agg
	}
	row := rows[0]
	v := fmt.Sprintf("%v", row[k]) //either this or a type assertion! :/
	byKey, present := agg[v]
	if !present {
		byKey = make([]Row, 0)
	}
	byKey = append(byKey, row)
	agg[v] = byKey
	return _groupByKey(k, rows[1:], agg)
}

func (rows Rows) SumOfKey(k string) int64 {
	return _sumOfKey(k, rows, 0)
}

func _sumOfKey(k string, rows []Row, agg int64) int64 {
	if len(rows) == 0 {
		return agg
	}
	return _sumOfKey(k, rows[1:], agg+rows[0][k].(int64))
}

type RowCondition = func(Row) bool

func (byKey GroupedRows) FilterGroups(condition RowCondition) GroupedRows {
	m := make(GroupedRows)
	for k, rows := range byKey {
		m[k] = Rows(rows).FilterRows(condition)
	}
	return m
}

func (rows Rows) FilterRows(condition RowCondition) []Row {
	agg := make([]Row, 0)
	return _filter(rows, condition, agg)
}

func _filter(rows []Row, condition RowCondition, agg []Row) []Row {
	if len(rows) == 0 {
		return agg
	}
	if condition(rows[0]) {
		agg = append(agg, rows[0])
	}
	return _filter(rows[1:], condition, agg)
}
