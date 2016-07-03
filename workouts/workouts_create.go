package workouts

import (
	"net/http"

	"github.com/jwfriese/workouttrackerapi/workouts/repository"
)

func WorkoutsCreateHandler(repository repository.WorkoutRepository) http.Handler {
	return nil
}

type workoutCreateHandler struct {
	repository repository.WorkoutRepository
}

func (h *workoutCreateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

}
