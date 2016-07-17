package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	liftdatamodel "github.com/jwfriese/workouttrackerapi/lifts/datamodel"
	setdatamodel "github.com/jwfriese/workouttrackerapi/lifts/sets/datamodel"
	setrepository "github.com/jwfriese/workouttrackerapi/lifts/sets/repository"
	"github.com/jwfriese/workouttrackerapi/sqlhelpers"
)

type LiftRepository interface {
	All() []*liftdatamodel.Lift
	GetById(id int) (*liftdatamodel.Lift, error)
	Insert(*liftdatamodel.Lift) (int, error)
}

type liftRepository struct {
	setRepository setrepository.SetRepository
	connection    *sql.DB
}

func NewLiftRepository(connection *sql.DB, repository setrepository.SetRepository) LiftRepository {
	return &liftRepository{
		setRepository: repository,
		connection:    connection,
	}
}

func (r *liftRepository) All() []*liftdatamodel.Lift {
	var lifts []*liftdatamodel.Lift
	rows, err := r.connection.Query("SELECT * from lifts;")

	if err != nil {
		log.Fatal(err)
	}

	var id int
	var name string
	var workout int
	var dataTemplate string
	var setIds sqlhelpers.IntSlice

	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&id, &name, &workout, &dataTemplate, &setIds)
		if err != nil {
			log.Fatal(err)
		}

		var sets []*setdatamodel.Set
		for _, setId := range setIds {
			set, _ := r.setRepository.GetById(setId)
			sets = append(sets, set)
		}

		lifts = append(lifts, &liftdatamodel.Lift{
			Id:           id,
			Name:         name,
			Workout:      workout,
			DataTemplate: dataTemplate,
			Sets:         sets,
		})
	}

	return lifts
}

func (r *liftRepository) GetById(id int) (*liftdatamodel.Lift, error) {
	query := fmt.Sprintf("SELECT * FROM lifts WHERE id=%v", id)
	row := r.connection.QueryRow(query)

	var liftId int
	var name string
	var workout int
	var dataTemplate string
	var setIds sqlhelpers.IntSlice

	err := row.Scan(&liftId, &name, &workout, &dataTemplate, &setIds)

	if err == sql.ErrNoRows {
		noLiftErrString := fmt.Sprintf("Lift with id=%v does not exist", id)
		return nil, errors.New(noLiftErrString)
	}

	if err != nil {
		return nil, err
	}

	var sets []*setdatamodel.Set
	for _, setId := range setIds {
		set, setErr := r.setRepository.GetById(setId)
		if setErr != nil {
			setErrString := fmt.Sprintf("Error fetching lift (id=%v): %s", id, setErr.Error())
			return nil, errors.New(setErrString)
		}

		sets = append(sets, set)
	}

	lift := &liftdatamodel.Lift{
		Id:           liftId,
		Name:         name,
		Workout:      workout,
		DataTemplate: dataTemplate,
		Sets:         sets,
	}

	return lift, nil
}

func (r *liftRepository) Insert(newLift *liftdatamodel.Lift) (int, error) {
	insertQuery := fmt.Sprintf("INSERT INTO lifts (name,workout,data_template) VALUES ('%v',%v,'%v') RETURNING id", newLift.Name, newLift.Workout, newLift.DataTemplate)
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

	var setIds sqlhelpers.IntSlice
	for _, set := range newLift.Sets {
		set.Lift = createdId
		setId, setInsertErr := r.setRepository.Insert(set)
		if setInsertErr != nil {
			setInsertErrString := fmt.Sprintf("Error inserting lift: %s", setInsertErr.Error())
			return -1, errors.New(setInsertErrString)
		}

		setIds = append(setIds, setId)
	}

	updateLiftQuery := fmt.Sprintf("UPDATE lifts SET sets='%s' WHERE id=%v", setIds.ToString(), createdId)
	_, updateLiftErr := r.connection.Exec(updateLiftQuery)

	if updateLiftErr != nil {
		return -1, updateLiftErr
	}

	return createdId, nil
}
