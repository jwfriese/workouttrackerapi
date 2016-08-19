package repository

import (
	"errors"
	"fmt"

	workoutdatamodel "github.com/jwfriese/workouttrackerapi/workouts/datamodel"
)

func (r *workoutRepository) Insert(workout *workoutdatamodel.Workout) (int, error) {
	if len(workout.Lifts) > 0 {
		errorString := fmt.Sprintf("Cannot include lift ids when inserting a workout (found %v)", workout.Lifts)
		return -1, errors.New(errorString)
	}

	insertStatement := fmt.Sprintf("INSERT INTO workouts (name,timestamp,lifts) VALUES ('%v','%v','{}') RETURNING id", workout.Name, workout.Timestamp)
	resultRows, err := r.connection.Query(insertStatement)
	if err != nil {
		return -1, err
	}

	defer resultRows.Close()
	resultRows.Next()

	var insertId int
	err = resultRows.Scan(&insertId)
	if err != nil {
		return -1, err
	}

	return insertId, nil
}
