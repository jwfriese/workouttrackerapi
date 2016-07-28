package repository

import (
	"errors"
	"fmt"

	liftdatamodel "github.com/jwfriese/workouttrackerapi/lifts/datamodel"
	"github.com/jwfriese/workouttrackerapi/sqlhelpers"
)

func (r *liftRepository) Insert(newLift *liftdatamodel.Lift) (int, error) {
	workoutExistsQuery := fmt.Sprintf("SELECT exists(SELECT 1 FROM workouts WHERE id=%v)", newLift.Workout)
	workoutExistsRow := r.connection.QueryRow(workoutExistsQuery)
	var workoutExists bool
	workoutScanErr := workoutExistsRow.Scan(&workoutExists)
	if workoutScanErr != nil {
		return -1, workoutScanErr
	}

	if !workoutExists {
		errString := fmt.Sprintf("Error inserting lift: Workout with id=%v does not exist", newLift.Workout)
		return -1, errors.New(errString)
	}

	insertQuery := fmt.Sprintf("INSERT INTO lifts (name,workout,data_template) VALUES ('%v',%v,'%v') RETURNING id", newLift.Name, newLift.Workout, newLift.DataTemplate)
	rows, err := r.connection.Query(insertQuery)
	if err != nil {
		return -1, err
	}

	defer rows.Close()
	rows.Next()

	var createdId int
	err = rows.Scan(&createdId)

	if err != nil {
		return -1, err
	}

	var setIds sqlhelpers.IntSlice
	for _, set := range newLift.Sets {
		set.Lift = createdId
		setId, setInsertErr := r.setRepository.Insert(set)
		if setInsertErr != nil {
			setInsertErrString := fmt.Sprintf("Error inserting lift: %s", setInsertErr.Error())
			return -1, errors.New(setInsertErrString)
		}

		setIds = append(setIds, setId)
	}

	updateLiftQuery := fmt.Sprintf("UPDATE lifts SET sets='%s' WHERE id=%v", setIds.ToString(), createdId)
	_, updateLiftErr := r.connection.Exec(updateLiftQuery)

	if updateLiftErr != nil {
		return -1, updateLiftErr
	}

	getWorkoutLiftsQuery := fmt.Sprintf("SELECT lifts FROM workouts WHERE id=%v", newLift.Workout)
	getWorkoutLiftsRow := r.connection.QueryRow(getWorkoutLiftsQuery)
	var liftIds sqlhelpers.IntSlice
	workoutLiftsScanErr := getWorkoutLiftsRow.Scan(&liftIds)
	if workoutLiftsScanErr != nil {
		return -1, workoutLiftsScanErr
	}

	liftIds = append(liftIds, createdId)

	updateWorkoutLiftsQuery := fmt.Sprintf("UPDATE workouts SET lifts='%v' WHERE id=%v", liftIds.ToString(), newLift.Workout)
	_, updateWorkoutLiftsErr := r.connection.Exec(updateWorkoutLiftsQuery)
	if updateWorkoutLiftsErr != nil {
		return -1, updateWorkoutLiftsErr
	}

	return createdId, nil
}
