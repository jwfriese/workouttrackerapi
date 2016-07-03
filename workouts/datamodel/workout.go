package datamodel

import (
	liftdatamodel "github.com/jwfriese/workouttrackerapi/lifts/datamodel"
)

type Workout struct {
	Id        int                   `json:"id"`
	Timestamp string                `json:"timestamp"`
	Lifts     []*liftdatamodel.Lift `json:"lifts"`
	Name      string                `json:"name"`
}
