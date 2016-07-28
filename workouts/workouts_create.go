package workouts

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
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
		writeRequestError(writer, readErr)
		return
	}
	request.Body.Close()

	workout, translateErr := handler.translator.Translate(requestBody)
	if translateErr != nil {
		writeRequestError(writer, translateErr)
		return
	}

	createdWorkoutId, insertErr := handler.repository.Insert(workout)
	if insertErr != nil {
		writeRequestError(writer, insertErr)
		return
	}

	createdWorkout, workoutFetchErr := handler.repository.GetById(createdWorkoutId)
	if workoutFetchErr != nil {
		writeServerError(writer, workoutFetchErr)
		return
	}

	createdWorkoutJSON, marshalErr := json.Marshal(&createdWorkout)
	if marshalErr != nil {
		writeServerError(writer, marshalErr)
		return
	}

	locationHeader := fmt.Sprintf("workouts/%v", createdWorkoutId)
	writer.Header().Set("Location", locationHeader)
	writer.WriteHeader(http.StatusCreated)
	writer.Write(createdWorkoutJSON)
}

func writeServerError(writer http.ResponseWriter, err error) {
	writer.WriteHeader(http.StatusInternalServerError)
	singleQuotedErr := bytes.Replace([]byte(err.Error()), []byte(`"`), []byte(`'`), -1)

	errString := fmt.Sprintf("{\"error\":\"%s\"}", string(singleQuotedErr))
	writer.Write([]byte(errString))
}

func writeRequestError(writer http.ResponseWriter, err error) {
	writer.WriteHeader(http.StatusBadRequest)
	singleQuotedErr := bytes.Replace([]byte(err.Error()), []byte(`"`), []byte(`'`), -1)

	errString := fmt.Sprintf("{\"error\":\"%s\"}", string(singleQuotedErr))
	writer.Write([]byte(errString))
}
