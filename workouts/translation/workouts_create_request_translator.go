package translation

import (
	"encoding/json"
	"errors"

	workoutdatamodel "github.com/jwfriese/workouttrackerapi/workouts/datamodel"
)

type WorkoutsCreateRequestTranslator interface {
	Translate(requestJSON []byte) (*workoutdatamodel.Workout, error)
}

type workoutsCreateRequestTranslator struct{}

func NewWorkoutsCreateRequestTranslator() WorkoutsCreateRequestTranslator {
	return &workoutsCreateRequestTranslator{}
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

	createdWorkout := &workoutdatamodel.Workout{
		Id:        -1,
		Name:      name,
		Timestamp: timestamp,
		Lifts:     nil,
	}

	return createdWorkout, nil
}

type workoutCreateRequest struct {
	Name      *string `json:"name"`
	Timestamp *string `json:"timestamp"`
}
