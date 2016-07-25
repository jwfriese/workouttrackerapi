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
	shouldAssociateWithWorkout := newLift.Workout > 0
	var liftIds sqlhelpers.IntSlice
	if shouldAssociateWithWorkout {
		associatedWorkoutQuery := fmt.Sprintf("SELECT lifts FROM workouts WHERE id=%v", newLift.Workout)
		associatedWorkoutRow := r.connection.QueryRow(associatedWorkoutQuery)

		workoutScanErr := associatedWorkoutRow.Scan(&liftIds)
		if workoutScanErr == sql.ErrNoRows {
			errorString := fmt.Sprintf("Error inserting lift: Workout with id=%v does not exist", newLift.Workout)
			return -1, errors.New(errorString)
		}

		if workoutScanErr != nil {
			return -1, workoutScanErr
		}
	}

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

	if shouldAssociateWithWorkout {
		liftIds = append(liftIds, createdId)
		updateWorkoutQuery := fmt.Sprintf("UPDATE workouts SET lifts='%s' WHERE id=%v", liftIds.ToString(), newLift.Workout)
		_, updateWorkoutErr := r.connection.Exec(updateWorkoutQuery)
		if updateWorkoutErr != nil {
			return -1, updateWorkoutErr
		}
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
