package translation

import (
	"encoding/json"
	"errors"

	liftdatamodel "github.com/jwfriese/workouttrackerapi/lifts/datamodel"
	lifttranslation "github.com/jwfriese/workouttrackerapi/lifts/translation"
	workoutdatamodel "github.com/jwfriese/workouttrackerapi/workouts/datamodel"
)

type WorkoutsCreateRequestTranslator interface {
	Translate(requestJSON []byte) (*workoutdatamodel.Workout, error)
}

type workoutsCreateRequestTranslator struct {
	liftTranslator lifttranslation.LiftsCreateRequestTranslator
}

func NewWorkoutsCreateRequestTranslator(liftTranslator lifttranslation.LiftsCreateRequestTranslator) WorkoutsCreateRequestTranslator {
	return &workoutsCreateRequestTranslator{
		liftTranslator: liftTranslator,
	}
}

func (translator *workoutsCreateRequestTranslator) Translate(requestJSON []byte) (*workoutdatamodel.Workout, error) {
	var workoutRequest workoutCreateRequest
	err := json.Unmarshal(requestJSON, &workoutRequest)
	if err != nil {
		return nil, err
	}

	if workoutRequest.Name == nil {
		return nil, errors.New("Missing required 'name' field from workout JSON")
	}

	name := *(workoutRequest.Name)

	if workoutRequest.Timestamp == nil {
		return nil, errors.New("Missing required 'timestamp' field from workout JSON")
	}

	timestamp := *(workoutRequest.Timestamp)

	if workoutRequest.LiftJSONObjects == nil {
		return nil, errors.New("Missing required 'lifts' field from workout JSON")
	}

	var lifts []*liftdatamodel.Lift
	for _, liftJSONInterface := range *(workoutRequest.LiftJSONObjects) {
		liftJSON, jsonErr := json.Marshal(liftJSONInterface)
		if jsonErr != nil {
			return nil, jsonErr
		}

		lift, liftErr := translator.liftTranslator.Translate(liftJSON)
		if liftErr != nil {
			return nil, liftErr
		}

		lifts = append(lifts, lift)
	}

	createdWorkout := &workoutdatamodel.Workout{
		Id:        -1,
		Name:      name,
		Timestamp: timestamp,
		Lifts:     lifts,
	}
	return createdWorkout, nil
}

type workoutCreateRequest struct {
	Name            *string        `json:"name"`
	Timestamp       *string        `json:"timestamp"`
	LiftJSONObjects *[]interface{} `json:"lifts"`
}
