package repository

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/jwfriese/workouttrackerapi/lifts/sets/datamodel"
	"github.com/jwfriese/workouttrackerapi/sqlhelpers"
)

type SetRepository interface {
	GetById(id int) *datamodel.Set
	Insert(set *datamodel.Set) (int, error)
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

	var (
		setId                 int
		dataTemplate          string
		lift                  int
		nullableWeight        sql.NullFloat64
		nullableHeight        sql.NullFloat64
		nullableTimeInSeconds sql.NullFloat64
		nullableReps          sql.NullInt64
	)

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

func (r *setRepository) Insert(set *datamodel.Set) (int, error) {
	weightString := sqlhelpers.Float32PointerToSQLString(set.Weight)
	heightString := sqlhelpers.Float32PointerToSQLString(set.Height)
	timeInSecondsString := sqlhelpers.Float32PointerToSQLString(set.TimeInSeconds)
	repsString := sqlhelpers.IntPointerToSQLString(set.Reps)

	insertQuery := fmt.Sprintf("INSERT INTO sets (data_template,lift,weight,height,time_in_seconds,reps) VALUES ('%v',%v,%v,%v,%v,%v) RETURNING id", set.DataTemplate, set.Lift, weightString, heightString, timeInSecondsString, repsString)

	rows, err := r.connection.Query(insertQuery)
	if err != nil {
		return -1, err
	}

	defer rows.Close()
	rows.Next()

	var createdId int
	err = rows.Scan(&createdId)
	if err != nil {
		return -1, err
	}

	return createdId, nil
}
