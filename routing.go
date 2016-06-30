package workouttrackerapi

import (
	"database/sql"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jwfriese/workouttrackerapi/lifts"
	liftrepository "github.com/jwfriese/workouttrackerapi/lifts/repository"
	setrepository "github.com/jwfriese/workouttrackerapi/lifts/sets/repository"
	"github.com/jwfriese/workouttrackerapi/workouts"
	workoutrepository "github.com/jwfriese/workouttrackerapi/workouts/repository"
)

func ApplicationHandler(db *sql.DB) http.Handler {
	handler := mux.NewRouter()

	workoutRepository := workoutrepository.NewWorkoutRepository(db)
	handler.Handle("/workouts", workouts.WorkoutsHandler(workoutRepository))

	setRepository := setrepository.NewSetRepository(db)
	liftRepository := liftrepository.NewLiftRepository(db, setRepository)
	handler.Handle("/lifts", lifts.LiftsHandler(liftRepository))

	return handler
}
