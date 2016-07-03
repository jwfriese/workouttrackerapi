package repository

import (
	"database/sql"
	"fmt"
	"log"

	liftdatamodel "github.com/jwfriese/workouttrackerapi/lifts/datamodel"
	setdatamodel "github.com/jwfriese/workouttrackerapi/lifts/sets/datamodel"
	"github.com/jwfriese/workouttrackerapi/lifts/sets/repository"
	"github.com/jwfriese/workouttrackerapi/sqlhelpers"
)

type LiftRepository interface {
	All() []*liftdatamodel.Lift
	GetById(id int) *liftdatamodel.Lift
}

type liftRepository struct {
	setRepository repository.SetRepository
	connection    *sql.DB
}

func NewLiftRepository(connection *sql.DB, repository repository.SetRepository) LiftRepository {
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
