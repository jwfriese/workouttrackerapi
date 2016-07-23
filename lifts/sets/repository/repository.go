package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/jwfriese/workouttrackerapi/lifts/sets/datamodel"
	"github.com/jwfriese/workouttrackerapi/sqlhelpers"
)

var ErrDoesNotExist = errors.New("Set with given id does not exist")

type SetRepository interface {
	GetById(id int) (*datamodel.Set, error)
	Insert(set *datamodel.Set) (int, error)
	Delete(id int) error
}

type setRepository struct {
	connection *sql.DB
}

func NewSetRepository(connection *sql.DB) SetRepository {
	return &setRepository{
		connection: connection,
	}
}

func (r *setRepository) GetById(id int) (*datamodel.Set, error) {
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

	err := row.Scan(&setId, &dataTemplate, &lift, &nullableWeight, &nullableHeight, &nullableTimeInSeconds, &nullableReps)

	if err == sql.ErrNoRows {
		noSetErrString := fmt.Sprintf("Set with id=%v does not exist", id)
		return nil, errors.New(noSetErrString)
	}

	if err != nil {
		return nil, err
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

	set := &datamodel.Set{
		Id:            setId,
		DataTemplate:  dataTemplate,
		Lift:          lift,
		Weight:        weight,
		Height:        height,
		TimeInSeconds: timeInSeconds,
		Reps:          reps,
	}

	return set, nil
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

func (repository *setRepository) Delete(id int) error {
	setQuery := fmt.Sprintf("SELECT id FROM sets WHERE id=%v", id)
	setRow := repository.connection.QueryRow(setQuery)

	var setId int
	scanErr := setRow.Scan(&setId)

	if scanErr == sql.ErrNoRows {
		return ErrDoesNotExist
	}

	deleteSetQuery := fmt.Sprintf("DELETE FROM sets WHERE id=%v", id)
	_, deleteSetErr := repository.connection.Exec(deleteSetQuery)

	return deleteSetErr
}
