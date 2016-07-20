package workouttrackerapi

import (
	"database/sql"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jwfriese/workouttrackerapi/lifts"
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
	handler.Handle("/lifts", lifts.LiftsHandler(liftRepository))

	workoutRepository := workoutrepository.NewWorkoutRepository(db, liftRepository)
	handler.Handle(workouts.WorkoutsShowEndpoint(), workouts.WorkoutsShowHandler(workoutRepository))

	setsCreateRequestTranslator := setstranslation.NewSetsCreateRequestTranslator()
	liftsCreateRequestTranslator := liftstranslation.NewLiftsCreateRequestTranslator(setsCreateRequestTranslator)
	workoutsCreateRequestTranslator := workoutstranslation.NewWorkoutsCreateRequestTranslator(liftsCreateRequestTranslator)

	workoutsEndpointHandler := workouts.WorkoutsEndpointHandler(workoutRepository, workoutsCreateRequestTranslator)
	handler.Handle("/workouts", workoutsEndpointHandler).Methods("GET", "POST")

	return handler
}
