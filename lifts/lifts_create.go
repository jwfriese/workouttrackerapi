package lifts

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jwfriese/workouttrackerapi/lifts/repository"
	"github.com/jwfriese/workouttrackerapi/lifts/translation"
)

func LiftsCreateHandler(repository repository.LiftRepository, translator translation.LiftsCreateRequestTranslator) http.Handler {
	return &liftsCreateHandler{
		repository: repository,
		translator: translator,
	}
}

type liftsCreateHandler struct {
	repository repository.LiftRepository
	translator translation.LiftsCreateRequestTranslator
}

func (handler *liftsCreateHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	queryVars := mux.Vars(request)
	workoutIdString := queryVars["workoutId"]
	isValidWorkoutParam, _ := regexp.Match("[1-9]+", []byte(workoutIdString))
	if !isValidWorkoutParam {
		err := errors.New("Path component following '/workouts' must be a valid workout id")
		writeErrorResponse(writer, err, http.StatusNotFound)
		return
	}

	workoutId, conversionErr := strconv.Atoi(workoutIdString)
	if conversionErr != nil {
		writeErrorResponse(writer, conversionErr, http.StatusInternalServerError)
		return
	}

	liftRequestJSON, liftRequestReadErr := ioutil.ReadAll(request.Body)
	request.Body.Close()

	if liftRequestReadErr != nil {
		writeErrorResponse(writer, liftRequestReadErr, http.StatusBadRequest)
		return
	}

	liftRequestModel, translationErr := handler.translator.Translate(liftRequestJSON)
	if translationErr != nil {
		writeErrorResponse(writer, translationErr, http.StatusBadRequest)
		return
	}

	liftRequestModel.Workout = workoutId

	createdLiftId, insertErr := handler.repository.Insert(liftRequestModel)
	if insertErr != nil {
		writeErrorResponse(writer, insertErr, http.StatusNotFound)
		return
	}

	fetchedLift, fetchErr := handler.repository.GetById(createdLiftId)
	if fetchErr != nil {
		writeErrorResponse(writer, fetchErr, http.StatusInternalServerError)
		return
	}

	marshaledLiftJSON, marshalErr := json.Marshal(fetchedLift)
	if marshalErr != nil {
		writeErrorResponse(writer, marshalErr, http.StatusInternalServerError)
		return
	}

	locationHeader := fmt.Sprintf("lifts/%v", createdLiftId)
	writer.Header().Set("Location", locationHeader)
	writer.WriteHeader(http.StatusCreated)
	writer.Write(marshaledLiftJSON)
}

func writeErrorResponse(writer http.ResponseWriter, err error, statusCode int) {
	writer.WriteHeader(statusCode)

	errorJSON := fmt.Sprintf("{\"error\":\"%s\"}", err.Error())
	writer.Write([]byte(errorJSON))
}
