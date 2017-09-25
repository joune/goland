package geo

import (
	"github.com/golang/geo/s2"
	"github.com/joune/zenly/data"
)

type Place struct {
	Lat, Lng float64
}

//enum to represent identified locations
type Location int8

const (
	Other Location = iota
	Home
	Work
	SameHome
)

func (loc Location) String() string {
	if loc == Home {
		return "Home"
	} else if loc == SameHome {
		return "SameHome"
	} else if loc == Work {
		return "Work"
	} else {
		return "Other"
	}
}

func (place Place) CellId() uint64 {
	//FIXME ? Should I use the full CellId or just the Face?
	return uint64(s2.CellFromLatLng(s2.LatLngFromDegrees(place.Lat, place.Lng)).ID().Parent(16))
}

func IdentifyLocation(usr1, usr2 data.User, lat, lng float64) Location {
	cellId := Place{lat, lng}.CellId()
	if cellId == usr1.GetHomeCell() {
		if cellId == usr2.GetHomeCell() {
			return SameHome
		} else {
			return Home
		}
	} else if cellId == usr1.GetWorkCell() {
		return Work
	} else {
		return Other
	}
}
