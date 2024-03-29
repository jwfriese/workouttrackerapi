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

var ErrDoesNotExist = errors.New("Lift with given id does not exist")

type LiftRepository interface {
	All() []*liftdatamodel.Lift
	GetById(id int) (*liftdatamodel.Lift, error)
	Insert(*liftdatamodel.Lift) (int, error)
	Delete(id int) error
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

		sets := []*setdatamodel.Set{}
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

	sets := []*setdatamodel.Set{}
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

func (repository *liftRepository) Delete(id int) error {
	liftQuery := fmt.Sprintf("SELECT id, sets FROM lifts WHERE id=%v", id)
	liftRow := repository.connection.QueryRow(liftQuery)

	var liftId int
	var setIds sqlhelpers.IntSlice

	scanErr := liftRow.Scan(&liftId, &setIds)
	if scanErr == sql.ErrNoRows {
		return ErrDoesNotExist
	}

	for _, setId := range setIds {
		deleteSetErr := repository.setRepository.Delete(setId)
		if deleteSetErr != nil && deleteSetErr != setrepository.ErrDoesNotExist {
			errString := fmt.Sprintf("Lift delete failed: Could not delete set with id=%v", setId)
			return errors.New(errString)
		}
	}

	deleteLiftQuery := fmt.Sprintf("DELETE FROM lifts WHERE id=%v", liftId)
	_, deleteLiftErr := repository.connection.Exec(deleteLiftQuery)

	return deleteLiftErr
}
