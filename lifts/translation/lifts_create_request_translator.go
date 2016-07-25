package translation

import (
	"encoding/json"
	"errors"

	liftdatamodel "github.com/jwfriese/workouttrackerapi/lifts/datamodel"
	setdatamodel "github.com/jwfriese/workouttrackerapi/lifts/sets/datamodel"
	settranslation "github.com/jwfriese/workouttrackerapi/lifts/sets/translation"
)

type LiftsCreateRequestTranslator interface {
	Translate(liftJSON []byte) (*liftdatamodel.Lift, error)
}

type liftsCreateRequestTranslator struct {
	setTranslator settranslation.SetsCreateRequestTranslator
}

func NewLiftsCreateRequestTranslator(setTranslator settranslation.SetsCreateRequestTranslator) LiftsCreateRequestTranslator {
	return &liftsCreateRequestTranslator{
		setTranslator: setTranslator,
	}
}

type liftsCreateRequest struct {
	Name         *string        `json:"name"`
	DataTemplate *string        `json:"dataTemplate"`
	Workout      *int           `json:"workout"`
	SetObjects   *[]interface{} `json:"sets"`
}

func (translator *liftsCreateRequestTranslator) Translate(requestJSON []byte) (*liftdatamodel.Lift, error) {
	var request liftsCreateRequest
	err := json.Unmarshal(requestJSON, &request)
	if err != nil {
		return nil, err
	}

	if request.Name == nil {
		return nil, errors.New("Missing required 'name' field in lift request JSON")
	}

	name := *(request.Name)

	if request.DataTemplate == nil {
		return nil, errors.New("Missing required 'dataTemplate' field in lift request JSON")
	}

	dataTemplate := *(request.DataTemplate)

	if request.SetObjects == nil {
		return nil, errors.New("Missing required 'sets' field in lift request JSON")
	}

	sets := []*setdatamodel.Set{}
	for _, setObject := range *(request.SetObjects) {
		setJSON, jsonErr := json.Marshal(setObject)
		if jsonErr != nil {
			return nil, jsonErr
		}

		set, translateErr := translator.setTranslator.Translate(setJSON)

		if translateErr != nil {
			return nil, translateErr
		}

		sets = append(sets, set)
	}

	workout := -1
	if request.Workout != nil {
		workout = *(request.Workout)
	}

	lift := &liftdatamodel.Lift{
		Id:           -1,
		Name:         name,
		DataTemplate: dataTemplate,
		Workout:      workout,
		Sets:         sets,
	}

	return lift, nil
}
