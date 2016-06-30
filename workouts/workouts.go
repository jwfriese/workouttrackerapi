package workouts

import (
	"bytes"
	_ "database/sql"
	"encoding/json"
	"github.com/jwfriese/workouttrackerapi/workouts/repository"
	_ "github.com/lib/pq"
	"log"
	"net/http"
)

func WorkoutsHandler(r repository.WorkoutRepository) http.Handler {
	return &workoutEndpoint{
		repository: r,
	}
}

type workoutEndpoint struct {
	repository repository.WorkoutRepository
}

func (we *workoutEndpoint) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	workouts := we.repository.All()

	var bytes bytes.Buffer

	workoutJSON, err := json.Marshal(workouts)
	if err != nil {
		log.Fatal(err)
	}

	bytes.Write(workoutJSON)

	json := bytes.Bytes()
	w.Write(json)
}
