package lifts

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	liftrepository "github.com/jwfriese/workouttrackerapi/lifts/repository"
	workoutrepository "github.com/jwfriese/workouttrackerapi/workouts/repository"
)

func LiftsShowHandler(liftRepository liftrepository.LiftRepository, workoutRepository workoutrepository.WorkoutRepository) http.Handler {
	return &liftsShowHandler{
		liftRepository:    liftRepository,
		workoutRepository: workoutRepository,
	}
}

type liftsShowHandler struct {
	liftRepository    liftrepository.LiftRepository
	workoutRepository workoutrepository.WorkoutRepository
}

func (h *liftsShowHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	queryArgs := mux.Vars(request)
	workoutId, err := strconv.Atoi(queryArgs["workoutId"])
	if err != nil {
		log.Fatal(err)
	}

	liftId, err := strconv.Atoi(queryArgs["liftId"])
	if err != nil {
		log.Fatal(err)
	}

	workout, _ := h.workoutRepository.GetById(workoutId)
	if workout == nil {
		writer.WriteHeader(http.StatusNotFound)
		responseBody := fmt.Sprintf("{\"error\":\"Workout with id=%v does not exist\"}", workoutId)
		_, _ = writer.Write([]byte(responseBody))
		return
	}

	lift, err := h.liftRepository.GetById(liftId)
	if err != nil {
		writer.WriteHeader(http.StatusNotFound)
		responseBody := fmt.Sprintf("{\"error\":\"Lift with id=%v does not exist\"}", liftId)
		_, _ = writer.Write([]byte(responseBody))
		return
	}

	if lift.Workout != workoutId {
		writer.WriteHeader(http.StatusNotFound)
		responseBody := fmt.Sprintf("{\"error\":\"Lift with id=%v does not exist on workout with id=%v\"}", liftId, workoutId)
		_, _ = writer.Write([]byte(responseBody))
		return
	}

	json, err := json.Marshal(lift)
	if err != nil {
		log.Fatal(err)
	}

	writer.WriteHeader(http.StatusOK)
	_, err = writer.Write(json)
	if err != nil {
		log.Fatal(err)
	}
}
