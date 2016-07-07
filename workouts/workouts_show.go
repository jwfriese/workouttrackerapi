package workouts

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jwfriese/workouttrackerapi/workouts/repository"
)

func WorkoutsShowEndpoint() string {
	return "/workouts/{id:[0-9]+}"
}

func WorkoutsShowHandler(r repository.WorkoutRepository) http.Handler {
	return &workoutsShowHandler{
		repository: r,
	}
}

type workoutsShowHandler struct {
	repository repository.WorkoutRepository
}

func (h *workoutsShowHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	args := mux.Vars(r)
	id, err := strconv.Atoi(args["id"])

	if err != nil {
		log.Fatal(err)
	}

	workout := h.repository.GetById(id)
	if workout == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	workoutBytes, jsonErr := json.Marshal(workout)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(workoutBytes)
}
