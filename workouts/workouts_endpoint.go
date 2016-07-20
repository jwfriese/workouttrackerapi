package workouts

import (
	"net/http"

	"github.com/jwfriese/workouttrackerapi/workouts/repository"
	"github.com/jwfriese/workouttrackerapi/workouts/translation"
)

func WorkoutsEndpointHandler(repository repository.WorkoutRepository, translator translation.WorkoutsCreateRequestTranslator) http.Handler {
	return &workoutsEndpointHandler{
		workoutsIndexHandler:  WorkoutsIndexHandler(repository),
		workoutsCreateHandler: WorkoutsCreateHandler(repository, translator),
	}
}

type workoutsEndpointHandler struct {
	workoutsIndexHandler  http.Handler
	workoutsCreateHandler http.Handler
}

func (handler *workoutsEndpointHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if request.Method == "GET" {
		handler.workoutsIndexHandler.ServeHTTP(writer, request)
	} else if request.Method == "POST" {
		handler.workoutsCreateHandler.ServeHTTP(writer, request)
	}
}
