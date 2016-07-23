package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	liftdatamodel "github.com/jwfriese/workouttrackerapi/lifts/datamodel"
	liftrepository "github.com/jwfriese/workouttrackerapi/lifts/repository"
	"github.com/jwfriese/workouttrackerapi/sqlhelpers"
	workoutdatamodel "github.com/jwfriese/workouttrackerapi/workouts/datamodel"
	_ "github.com/lib/pq"
)

var ErrDoesNotExist = errors.New("Workout with the given id does not exist")

type WorkoutRepository interface {
	All() []*workoutdatamodel.Workout
	GetById(id int) (*workoutdatamodel.Workout, error)
	Insert(workout *workoutdatamodel.Workout) (int, error)
	Delete(id int) error
}

type workoutRepository struct {
	connection     *sql.DB
	liftRepository liftrepository.LiftRepository
}

func NewWorkoutRepository(db *sql.DB, liftRepository liftrepository.LiftRepository) WorkoutRepository {
	return &workoutRepository{
		connection:     db,
		liftRepository: liftRepository,
	}
}

func (r *workoutRepository) All() []*workoutdatamodel.Workout {
	if r.connection != nil {
		rows, err := r.connection.Query("SELECT * FROM workouts")
		if err != nil {
			log.Fatal(err)
		}

		workouts := []*workoutdatamodel.Workout{}

		var id int
		var name string
		var timestamp string
		var liftIds sqlhelpers.IntSlice

		defer rows.Close()
		for rows.Next() {
			err := rows.Scan(&id, &name, &timestamp, &liftIds)
			if err != nil {
				log.Fatal(err)
			}

			lifts := []*liftdatamodel.Lift{}
			for _, liftId := range liftIds {
				lift, _ := r.liftRepository.GetById(liftId)
				lifts = append(lifts, lift)
			}

			workouts = append(workouts, &workoutdatamodel.Workout{
				Id:        id,
				Name:      name,
				Timestamp: timestamp,
				Lifts:     lifts,
			})
		}

		return workouts
	}

	return nil
}

func (r *workoutRepository) GetById(id int) (*workoutdatamodel.Workout, error) {
	query := fmt.Sprintf("SELECT * FROM workouts WHERE id = %v", id)
	row := r.connection.QueryRow(query)

	var workoutId int
	var name string
	var timestamp string
	var liftIds sqlhelpers.IntSlice

	err := row.Scan(&workoutId, &name, &timestamp, &liftIds)
	if err == sql.ErrNoRows {
		noResultsErrString := fmt.Sprintf("Workout with id=%v does not exist", id)
		return nil, errors.New(noResultsErrString)
	}

	// How do I test handling of arbitrary errors?
	if err != nil {
		return nil, err
	}

	lifts := []*liftdatamodel.Lift{}
	for _, liftId := range liftIds {
		lift, liftErr := r.liftRepository.GetById(liftId)
		if liftErr != nil {
			liftErrString := fmt.Sprintf("Error fetching workout (id=%v): %s", id, liftErr.Error())
			return nil, errors.New(liftErrString)
		}

		lifts = append(lifts, lift)
	}

	workout := &workoutdatamodel.Workout{
		Id:        workoutId,
		Name:      name,
		Timestamp: timestamp,
		Lifts:     lifts,
	}

	return workout, nil
}

func (r *workoutRepository) Insert(workout *workoutdatamodel.Workout) (int, error) {
	insertStatement := fmt.Sprintf("INSERT INTO workouts (name,timestamp) VALUES ('%v','%v') RETURNING id", workout.Name, workout.Timestamp)
	resultRows, err := r.connection.Query(insertStatement)
	if err != nil {
		log.Fatal(err)
	}

	defer resultRows.Close()
	resultRows.Next()

	var insertId int
	err = resultRows.Scan(&insertId)
	if err != nil {
		return -1, err
	}

	var liftIds sqlhelpers.IntSlice
	for _, lift := range workout.Lifts {
		lift.Workout = insertId
		liftId, liftInsertErr := r.liftRepository.Insert(lift)
		if liftInsertErr != nil {
			liftInsertErrString := fmt.Sprintf("Error inserting workout: %s", liftInsertErr.Error())
			return -1, errors.New(liftInsertErrString)
		}

		liftIds = append(liftIds, liftId)
	}

	updateLiftIdsQuery := fmt.Sprintf("UPDATE workouts SET lifts='%v' WHERE id=%v", liftIds.ToString(), insertId)
	_, updateErr := r.connection.Exec(updateLiftIdsQuery)
	if updateErr != nil {
		return -1, updateErr
	}

	return insertId, nil
}

func (repository *workoutRepository) Delete(id int) error {
	workoutQuery := fmt.Sprintf("SELECT id, lifts FROM workouts WHERE id=%v", id)
	workoutRow := repository.connection.QueryRow(workoutQuery)

	var workoutId int
	var lifts sqlhelpers.IntSlice

	scanErr := workoutRow.Scan(&workoutId, &lifts)
	if scanErr == sql.ErrNoRows {
		return ErrDoesNotExist
	}

	for _, liftId := range lifts {
		deleteLiftErr := repository.liftRepository.Delete(liftId)
		if deleteLiftErr != nil && deleteLiftErr != liftrepository.ErrDoesNotExist {
			errString := fmt.Sprintf("Workout failed to delete: Lift with id=%v could not be deleted", liftId)
			return errors.New(errString)
		}
	}

	deleteWorkoutQuery := fmt.Sprintf("DELETE FROM workouts WHERE id=%v", workoutId)
	_, deleteErr := repository.connection.Exec(deleteWorkoutQuery)

	// How do I test this deleteErr?
	return deleteErr
}
