package repository

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/jwfriese/workouttrackerapi/lifts/sets/datamodel"
)

type SetRepository interface {
	GetById(id int) *datamodel.Set
}

type setRepository struct {
	connection *sql.DB
}

func NewSetRepository(connection *sql.DB) SetRepository {
	return &setRepository{
		connection: connection,
	}
}

func (r *setRepository) GetById(id int) *datamodel.Set {
	queryString := fmt.Sprintf("SELECT * FROM sets WHERE id = %v", id)
	row := r.connection.QueryRow(queryString)

	var setId int
	var dataTemplate string
	var lift int
	var nullableWeight sql.NullFloat64
	var nullableHeight sql.NullFloat64
	var nullableTimeInSeconds sql.NullFloat64
	var nullableReps sql.NullInt64

	var err error
	err = row.Scan(&setId, &dataTemplate, &lift, &nullableWeight, &nullableHeight, &nullableTimeInSeconds, &nullableReps)

	if err != nil {
		log.Fatal(err)
	}

	var weight *float32
	if nullableWeight.Valid {
		weight = new(float32)
		(*weight) = float32(nullableWeight.Float64)
	}

	var height *float32
	if nullableHeight.Valid {
		height = new(float32)
		(*height) = float32(nullableHeight.Float64)
	}

	var timeInSeconds *float32
	if nullableTimeInSeconds.Valid {
		timeInSeconds = new(float32)
		(*timeInSeconds) = float32(nullableTimeInSeconds.Float64)
	}

	var reps *int
	if nullableReps.Valid {
		reps = new(int)
		(*reps) = int(nullableReps.Int64)
	}

	return &datamodel.Set{
		Id:            setId,
		DataTemplate:  dataTemplate,
		Lift:          lift,
		Weight:        weight,
		Height:        height,
		TimeInSeconds: timeInSeconds,
		Reps:          reps,
	}
}
