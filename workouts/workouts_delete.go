package workouts

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jwfriese/workouttrackerapi/workouts/repository"
)

func WorkoutsDeleteHandler(r repository.WorkoutRepository) http.Handler {
	return &workoutsDeleteHandler{
		repository: r,
	}
}

type workoutsDeleteHandler struct {
	repository repository.WorkoutRepository
}

func (handler *workoutsDeleteHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	toDeleteId, idReadErr := strconv.Atoi(vars["id"])
	if idReadErr != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	deleteErr := handler.repository.Delete(toDeleteId)
	if deleteErr != nil {
		writer.WriteHeader(http.StatusNotFound)
	} else {
		writer.WriteHeader(http.StatusNoContent)
	}
}
