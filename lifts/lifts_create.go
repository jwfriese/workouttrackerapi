package lifts

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

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
	requestBodyBytes, bodyReadErr := ioutil.ReadAll(request.Body)
	request.Body.Close()
	if bodyReadErr != nil {
		log.Fatal(bodyReadErr)
	}

	newLiftModel, translationErr := handler.translator.Translate(requestBodyBytes)
	if translationErr != nil {
		log.Fatal(translationErr)
	}

	createdLiftId, insertErr := handler.repository.Insert(newLiftModel)
	if insertErr != nil {
		log.Fatal(insertErr)
	}

	createdLift, liftFetchErr := handler.repository.GetById(createdLiftId)
	if liftFetchErr != nil {
		log.Fatal(liftFetchErr)
	}

	createdLiftJSON, jsonMarshalErr := json.Marshal(createdLift)
	if jsonMarshalErr != nil {
		log.Fatal(jsonMarshalErr)
	}

	locationHeader := fmt.Sprintf("/lifts/%v", createdLiftId)
	writer.Header().Add("Location", locationHeader)
	writer.WriteHeader(http.StatusCreated)
	writer.Write(createdLiftJSON)
}
