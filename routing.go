package workouttracker

import (
	"github.com/gorilla/mux"
	"github.com/jwfriese/workouttracker/workouts"
	"github.com/jwfriese/workouttracker/workouts/repository"

	"database/sql"
	"net/http"
)

func ApplicationHandler(db *sql.DB) http.Handler {
	handler := mux.NewRouter()

	workoutRepository := repository.NewWorkoutRepository(db)
	handler.Handle("/workouts", workouts.WorkoutsHandler(workoutRepository))
	//handler.HandleFunc("/workouts/{id:[0-9]+}", workouts.WorkoutHandler)

	return handler
}
