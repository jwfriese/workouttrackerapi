package datamodel

import (
	"github.com/jwfriese/workouttrackerapi/lifts/sets/datamodel"
)

type Lift struct {
	Id           int
	Name         string
	DataTemplate string
	Workout      int
	Sets         []*datamodel.Set
}
