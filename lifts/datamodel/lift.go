package datamodel

import (
	"github.com/jwfriese/workouttrackerapi/lifts/sets/datamodel"
)

type Lift struct {
	Id           int              `json:"id"`
	Name         string           `json:"name"`
	DataTemplate string           `json:"dataTemplate"`
	Workout      int              `json:"workout"`
	Sets         []*datamodel.Set `json:"sets"`
}
