package datamodel

import (
	liftdatamodel "github.com/jwfriese/workouttrackerapi/lifts/datamodel"
)

type Workout struct {
	Id        int
	Timestamp string
	Lifts     []*liftdatamodel.Lift
	Name      string
}
