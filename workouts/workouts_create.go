package workouts

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/jwfriese/workouttrackerapi/workouts/repository"
	"github.com/jwfriese/workouttrackerapi/workouts/translation"
)

func WorkoutsCreateHandler(r repository.WorkoutRepository, t translation.WorkoutsCreateRequestTranslator) http.Handler {
	return &workoutsCreateHandler{
		repository: r,
		translator: t,
	}
}

type workoutsCreateHandler struct {
	repository repository.WorkoutRepository
	translator translation.WorkoutsCreateRequestTranslator
}

func (handler *workoutsCreateHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "application/json;charset=UTF-8")

	requestBody, readErr := ioutil.ReadAll(request.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}
	request.Body.Close()

	workout, translateErr := handler.translator.Translate(requestBody)
	if translateErr != nil {
		log.Fatal(translateErr)
	}
	createdWorkoutId, insertErr := handler.repository.Insert(workout)
	if insertErr != nil {
		log.Fatal(insertErr)
	}

	locationHeader := fmt.Sprintf("workouts/%v", createdWorkoutId)
	writer.Header().Set("Location", locationHeader)

	writer.WriteHeader(http.StatusCreated)

	createdWorkout := handler.repository.GetById(createdWorkoutId)
	createdWorkoutJSON, marshalErr := json.Marshal(&createdWorkout)
	if marshalErr != nil {
		log.Fatal(marshalErr)
	}

	writer.Write(createdWorkoutJSON)
}
