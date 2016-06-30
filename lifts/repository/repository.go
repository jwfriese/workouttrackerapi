package repository

import (
	"database/sql"
	"log"

	liftdatamodel "github.com/jwfriese/workouttrackerapi/lifts/datamodel"
	setdatamodel "github.com/jwfriese/workouttrackerapi/lifts/sets/datamodel"
	"github.com/jwfriese/workouttrackerapi/lifts/sets/repository"
	"github.com/jwfriese/workouttrackerapi/sqlhelpers"
)

type LiftRepository interface {
	All() []*liftdatamodel.Lift
	GetById(id string) *liftdatamodel.Lift
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
	var setIds sqlhelpers.UIntSlice

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

func (r *liftRepository) GetById(id string) *liftdatamodel.Lift {
	return nil
}
