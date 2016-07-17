package translation

import (
	"encoding/json"
	"errors"
	"fmt"

	setdatamodel "github.com/jwfriese/workouttrackerapi/lifts/sets/datamodel"
)

type SetsCreateRequestTranslator interface {
	Translate(requestJSON []byte) (*setdatamodel.Set, error)
}

type setsCreateRequestTranslator struct{}

func NewSetsCreateRequestTranslator() SetsCreateRequestTranslator {
	return &setsCreateRequestTranslator{}
}

type setsCreateRequest struct {
	DataTemplate  *string  `json:"dataTemplate"`
	Weight        *float32 `json:"weight"`
	Height        *float32 `json:"height"`
	TimeInSeconds *float32 `json:"timeInSeconds"`
	Reps          *int     `json:"reps"`
}

func (translator *setsCreateRequestTranslator) Translate(requestJSON []byte) (*setdatamodel.Set, error) {
	var setRequest *setsCreateRequest
	requestErr := json.Unmarshal(requestJSON, &setRequest)
	if requestErr != nil {
		return nil, requestErr
	}

	validationErr := validateRequest(setRequest)
	if validationErr != nil {
		return nil, validationErr
	}

	set := createSetFromRequest(setRequest)
	return set, nil
}

func validateRequest(request *setsCreateRequest) error {
	if request.DataTemplate == nil {
		return errors.New("Missing required \"dataTemplate\" field in set request")
	}

	if *(request.DataTemplate) == "weight/reps" {
		return validateWeightRepsRequest(request)
	} else if *(request.DataTemplate) == "height/reps" {
		return validateHeightRepsRequest(request)
	} else if *(request.DataTemplate) == "timeInSeconds" {
		return validateTimeInSecondsRequest(request)
	} else if *(request.DataTemplate) == "weight/timeInSeconds" {
		return validateWeightTimeInSecondsRequest(request)
	}

	errString := fmt.Sprintf("Unrecognized data template \"%s\"", *(request.DataTemplate))
	return errors.New(errString)
}

func validateWeightRepsRequest(request *setsCreateRequest) error {
	if request.Weight == nil {
		return errors.New("Missing required \"weight\" field in request for \"weight/reps\" set")
	}

	if request.Reps == nil {
		return errors.New("Missing required \"reps\" field in request for \"weight/reps\" set")
	}

	return nil
}

func validateHeightRepsRequest(request *setsCreateRequest) error {
	if request.Height == nil {
		return errors.New("Missing required \"height\" field in request for \"height/reps\" set")
	}

	if request.Reps == nil {
		return errors.New("Missing required \"reps\" field in request for \"height/reps\" set")
	}

	return nil
}

func validateTimeInSecondsRequest(request *setsCreateRequest) error {
	if request.TimeInSeconds == nil {
		return errors.New("Missing required \"timeInSeconds\" field in request for \"timeInSeconds\" set")
	}

	return nil
}

func validateWeightTimeInSecondsRequest(request *setsCreateRequest) error {
	if request.Weight == nil {
		return errors.New("Missing required \"weight\" field in request for \"weight/timeInSeconds\" set")
	}

	if request.TimeInSeconds == nil {
		return errors.New("Missing required \"timeInSeconds\" field in request for \"weight/timeInSeconds\" set")
	}

	return nil
}

func createSetFromRequest(request *setsCreateRequest) *setdatamodel.Set {
	if *(request.DataTemplate) == "weight/reps" {
		return createWeightRepsSet(request)
	} else if *(request.DataTemplate) == "height/reps" {
		return createHeightRepsSet(request)
	} else if *(request.DataTemplate) == "timeInSeconds" {
		return createTimeInSecondsSet(request)
	} else if *(request.DataTemplate) == "weight/timeInSeconds" {
		return createWeightTimeInSecondsSet(request)
	}

	return nil
}

func createWeightRepsSet(request *setsCreateRequest) *setdatamodel.Set {
	return &setdatamodel.Set{
		Id:           -1,
		DataTemplate: "weight/reps",
		Lift:         -1,
		Weight:       request.Weight,
		Reps:         request.Reps,
	}
}

func createHeightRepsSet(request *setsCreateRequest) *setdatamodel.Set {
	return &setdatamodel.Set{
		Id:           -1,
		DataTemplate: "height/reps",
		Lift:         -1,
		Height:       request.Height,
		Reps:         request.Reps,
	}
}

func createTimeInSecondsSet(request *setsCreateRequest) *setdatamodel.Set {
	return &setdatamodel.Set{
		Id:            -1,
		DataTemplate:  "timeInSeconds",
		Lift:          -1,
		TimeInSeconds: request.TimeInSeconds,
	}
}

func createWeightTimeInSecondsSet(request *setsCreateRequest) *setdatamodel.Set {
	return &setdatamodel.Set{
		Id:            -1,
		DataTemplate:  "weight/timeInSeconds",
		Lift:          -1,
		Weight:        request.Weight,
		TimeInSeconds: request.TimeInSeconds,
	}
}
