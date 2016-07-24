package workouttrackerapi

import (
	"database/sql"
	"net/http"

	"github.com/gorilla/mux"
	liftrepository "github.com/jwfriese/workouttrackerapi/lifts/repository"
	setrepository "github.com/jwfriese/workouttrackerapi/lifts/sets/repository"
	setstranslation "github.com/jwfriese/workouttrackerapi/lifts/sets/translation"
	liftstranslation "github.com/jwfriese/workouttrackerapi/lifts/translation"
	"github.com/jwfriese/workouttrackerapi/workouts"
	workoutrepository "github.com/jwfriese/workouttrackerapi/workouts/repository"
	workoutstranslation "github.com/jwfriese/workouttrackerapi/workouts/translation"
)

func ApplicationHandler(db *sql.DB) http.Handler {
	handler := mux.NewRouter()

	setRepository := setrepository.NewSetRepository(db)
	liftRepository := liftrepository.NewLiftRepository(db, setRepository)

	workoutRepository := workoutrepository.NewWorkoutRepository(db, liftRepository)
	handler.Handle("/workouts/{id:[0-9]+}", workouts.WorkoutsShowHandler(workoutRepository)).Methods("GET")
	handler.Handle("/workouts/{id}", workouts.WorkoutsDeleteHandler(workoutRepository)).Methods("DELETE")

	setsCreateRequestTranslator := setstranslation.NewSetsCreateRequestTranslator()
	liftsCreateRequestTranslator := liftstranslation.NewLiftsCreateRequestTranslator(setsCreateRequestTranslator)
	workoutsCreateRequestTranslator := workoutstranslation.NewWorkoutsCreateRequestTranslator(liftsCreateRequestTranslator)

	handler.Handle("/workouts", workouts.WorkoutsIndexHandler(workoutRepository)).Methods("GET")
	handler.Handle("/workouts", workouts.WorkoutsCreateHandler(workoutRepository, workoutsCreateRequestTranslator)).Methods("POST")

	return handler
}
