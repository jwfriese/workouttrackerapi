package repository

import (
	"database/sql"
	"fmt"
	"github.com/jwfriese/workouttracker/sqlhelpers"
	"github.com/jwfriese/workouttracker/workouts/datamodel"
	"log"
)

type WorkoutRepository interface {
	All() []*datamodel.Workout
	GetById(id string) *datamodel.Workout
}

type workoutRepository struct {
	connection *sql.DB
}

func NewWorkoutRepository(db *sql.DB) WorkoutRepository {
	return &workoutRepository{
		connection: db,
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
		var lifts sqlhelpers.SQLStringSlice
		defer rows.Close()
		for rows.Next() {
			err := rows.Scan(&id, &name, &timestamp, &lifts)
			if err != nil {
				log.Fatal(err)
			}

			workouts = append(workouts, &datamodel.Workout{
				Id:        id,
				Name:      name,
				Timestamp: timestamp,
				Lifts:     lifts.ToStringSlice(),
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
		var lifts sqlhelpers.SQLStringSlice
		defer rows.Close()
		for rows.Next() {
			err := rows.Scan(&id, &name, &timestamp, &lifts)
			if err != nil {
				log.Fatal(err)
			}

			return &datamodel.Workout{
				Id:        id,
				Name:      name,
				Timestamp: timestamp,
				Lifts:     lifts.ToStringSlice(),
			}
		}
	}

	return nil
}
