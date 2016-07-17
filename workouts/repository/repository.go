package repository

import (
	"database/sql"
	"fmt"
	"log"

	liftdatamodel "github.com/jwfriese/workouttrackerapi/lifts/datamodel"
	liftrepository "github.com/jwfriese/workouttrackerapi/lifts/repository"
	"github.com/jwfriese/workouttrackerapi/sqlhelpers"
	workoutdatamodel "github.com/jwfriese/workouttrackerapi/workouts/datamodel"
	_ "github.com/lib/pq"
)

type WorkoutRepository interface {
	All() []*workoutdatamodel.Workout
	GetById(id int) *workoutdatamodel.Workout
	Insert(workout *workoutdatamodel.Workout) (int, error)
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
				lift := r.liftRepository.GetById(liftId)
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

func (r *workoutRepository) GetById(id int) *workoutdatamodel.Workout {
	if r.connection != nil {
		query := fmt.Sprintf("SELECT * FROM workouts WHERE id = %v", id)
		rows, err := r.connection.Query(query)

		if err != nil {
			log.Fatal(err)
		}

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
				lift := r.liftRepository.GetById(liftId)
				lifts = append(lifts, lift)
			}

			return &workoutdatamodel.Workout{
				Id:        id,
				Name:      name,
				Timestamp: timestamp,
				Lifts:     lifts,
			}
		}
	}

	return nil
}

func (r *workoutRepository) Insert(workout *workoutdatamodel.Workout) (int, error) {
	var liftIds sqlhelpers.IntSlice
	for _, lift := range workout.Lifts {
		liftId, _ := r.liftRepository.Insert(lift)
		liftIds = append(liftIds, liftId)
	}

	insertStatement := fmt.Sprintf("INSERT INTO workouts (name,timestamp,lifts) VALUES ('%v','%v','%v') RETURNING id", workout.Name, workout.Timestamp, liftIds.ToString())
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

	return insertId, nil
}
