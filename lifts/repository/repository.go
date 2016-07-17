package repository

import (
	"database/sql"
	"fmt"
	"log"

	liftdatamodel "github.com/jwfriese/workouttrackerapi/lifts/datamodel"
	setdatamodel "github.com/jwfriese/workouttrackerapi/lifts/sets/datamodel"
	setrepository "github.com/jwfriese/workouttrackerapi/lifts/sets/repository"
	"github.com/jwfriese/workouttrackerapi/sqlhelpers"
)

type LiftRepository interface {
	All() []*liftdatamodel.Lift
	GetById(id int) *liftdatamodel.Lift
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
			set := r.setRepository.GetById(setId)
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

func (r *liftRepository) GetById(id int) *liftdatamodel.Lift {
	query := fmt.Sprintf("SELECT * FROM lifts WHERE id=%v", id)
	row := r.connection.QueryRow(query)

	var liftId int
	var name string
	var workout int
	var dataTemplate string
	var setIds sqlhelpers.IntSlice

	err := row.Scan(&liftId, &name, &workout, &dataTemplate, &setIds)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}

		log.Fatal(err)
	}

	var sets []*setdatamodel.Set
	for _, setId := range setIds {
		set := r.setRepository.GetById(setId)
		sets = append(sets, set)
	}

	return &liftdatamodel.Lift{
		Id:           liftId,
		Name:         name,
		Workout:      workout,
		DataTemplate: dataTemplate,
		Sets:         sets,
	}
}

func (r *liftRepository) Insert(newLift *liftdatamodel.Lift) (int, error) {
	var setIds sqlhelpers.IntSlice
	for _, set := range newLift.Sets {
		setId, _ := r.setRepository.Insert(set)
		setIds = append(setIds, setId)
	}

	insertQuery := fmt.Sprintf("INSERT INTO lifts (name,workout,data_template,sets) VALUES ('%v',%v,'%v','%v') RETURNING id", newLift.Name, newLift.Workout, newLift.DataTemplate, setIds.ToString())
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
