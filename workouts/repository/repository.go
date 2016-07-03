package repository

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/jwfriese/workouttrackerapi/lifts/repository"
	"github.com/jwfriese/workouttrackerapi/sqlhelpers"
	"github.com/jwfriese/workouttrackerapi/workouts/datamodel"
)

type WorkoutRepository interface {
	All() []*datamodel.Workout
	GetById(id string) *datamodel.Workout
}

type workoutRepository struct {
	connection     *sql.DB
	liftRepository repository.LiftRepository
}

func NewWorkoutRepository(db *sql.DB, liftRepository repository.LiftRepository) WorkoutRepository {
	return &workoutRepository{
		connection:     db,
		liftRepository: liftRepository,
	}
}

func (r *workoutRepository) All() []*datamodel.Workout {
	if r.connection != nil {
		rows, err := r.connection.Query("SELECT * FROM workouts")
		if err != nil {
			log.Fatal(err)
		}

		workouts := []*datamodel.Workout{}

		var id int
		var name string
		var timestamp string
		var liftIds sqlhelpers.UIntSlice

		defer rows.Close()
		for rows.Next() {
			err := rows.Scan(&id, &name, &timestamp, &liftIds)
			if err != nil {
				log.Fatal(err)
			}

			workouts = append(workouts, &datamodel.Workout{
				Id:        id,
				Name:      name,
				Timestamp: timestamp,
				Lifts:     liftIds,
			})
		}

		return workouts
	}

	return nil
}

func (r *workoutRepository) GetById(id string) *datamodel.Workout {
	if r.connection != nil {
		query := fmt.Sprintf("SELECT * FROM workouts WHERE id = %s", id)
		rows, err := r.connection.Query(query)

		if err != nil {
			log.Fatal(err)
		}

		var id int
		var name string
		var timestamp string
		var liftIds sqlhelpers.UIntSlice

		defer rows.Close()
		for rows.Next() {
			err := rows.Scan(&id, &name, &timestamp, &liftIds)
			if err != nil {
				log.Fatal(err)
			}

			return &datamodel.Workout{
				Id:        id,
				Name:      name,
				Timestamp: timestamp,
				Lifts:     liftIds,
			}
		}
	}

	return nil
}
