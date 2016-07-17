package repository_test

import (
	"database/sql"
	"fmt"
	"log"

	liftdatamodel "github.com/jwfriese/workouttrackerapi/lifts/datamodel"
	"github.com/jwfriese/workouttrackerapi/lifts/repository"
	setdatamodel "github.com/jwfriese/workouttrackerapi/lifts/sets/datamodel"
	"github.com/jwfriese/workouttrackerapi/lifts/sets/repository/repositoryfakes"
	"github.com/jwfriese/workouttrackerapi/sqlhelpers"
	_ "github.com/lib/pq"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func populateLiftsDatabase(openConnection *sql.DB) {
	_, err := openConnection.Exec("INSERT INTO lifts (name,workout,data_template,sets) VALUES ('turtle lift',1,'weight/reps','{1,2}')")
	if err != nil {
		log.Fatal(err)
	}

	_, err = openConnection.Exec("INSERT INTO lifts (name,workout,data_template,sets) VALUES ('turtle jump',3,'height/reps','{3,4}')")
	if err != nil {
		log.Fatal(err)
	}
}

func clearLiftsDatabase(openConnection *sql.DB) {
	_, err := openConnection.Exec("TRUNCATE lifts")
	if err != nil {
		log.Fatal(err)
	}

	_, err = openConnection.Exec("ALTER SEQUENCE lifts_id_seq RESTART WITH 1")
	if err != nil {
		log.Fatal(err)
	}
}

var _ = Describe("LiftRepository", func() {
	var (
		subject           repository.LiftRepository
		testConnection    *sql.DB
		fakeSetRepository *repositoryfakes.FakeSetRepository
	)

	BeforeEach(func() {
		var err error
		testConnection, err = sql.Open("postgres", "user=postgres dbname=workout_tracker_unit_test sslmode=disable")
		if err != nil {
			log.Fatal(err)
		}

		populateLiftsDatabase(testConnection)

		fakeSetRepository = new(repositoryfakes.FakeSetRepository)
		subject = repository.NewLiftRepository(testConnection, fakeSetRepository)
	})

	AfterEach(func() {
		clearLiftsDatabase(testConnection)
	})

	Describe("Getting lifts from the database", func() {
		var (
			turtleSet *setdatamodel.Set
			crabSet   *setdatamodel.Set
			puppySet  *setdatamodel.Set
			kittenSet *setdatamodel.Set
		)

		BeforeEach(func() {
			turtleSet = &setdatamodel.Set{}
			crabSet = &setdatamodel.Set{}
			puppySet = &setdatamodel.Set{}
			kittenSet = &setdatamodel.Set{}

			fakeSetRepository.GetByIdStub = func(id int) *setdatamodel.Set {
				if id == 1 {
					return turtleSet
				} else if id == 2 {
					return crabSet
				} else if id == 3 {
					return puppySet
				} else if id == 4 {
					return kittenSet
				}

				return nil
			}
		})

		Describe("All", func() {
			var (
				lifts []*liftdatamodel.Lift
			)

			BeforeEach(func() {
				lifts = subject.All()
			})

			It("should get all lifts from the database", func() {
				Expect(len(lifts)).To(Equal(2))
				Expect(lifts[0].Id).To(Equal(1))
				Expect(lifts[0].Name).To(Equal("turtle lift"))
				Expect(lifts[1].Id).To(Equal(2))
				Expect(lifts[1].Name).To(Equal("turtle jump"))

				Expect(len(lifts[0].Sets)).To(Equal(2))
				Expect(lifts[0].Sets[0]).To(BeIdenticalTo(turtleSet))
				Expect(lifts[0].Sets[1]).To(BeIdenticalTo(crabSet))
				Expect(lifts[1].Sets[0]).To(BeIdenticalTo(puppySet))
				Expect(lifts[1].Sets[1]).To(BeIdenticalTo(kittenSet))
			})
		})

		Describe("Getting a single lift by id", func() {
			Context("When there exists a lift with that ID in the database", func() {
				var (
					result *liftdatamodel.Lift
				)

				BeforeEach(func() {
					result = subject.GetById(2)
				})

				It("retrieves the requested lift from the database", func() {
					Expect(result).ToNot(BeNil())

					Expect(result.Id).To(Equal(2))
					Expect(result.Name).To(Equal("turtle jump"))
					Expect(result.Workout).To(Equal(3))
					Expect(result.DataTemplate).To(Equal("height/reps"))

					Expect(len(result.Sets)).To(Equal(2))
					Expect(result.Sets[0]).To(BeIdenticalTo(puppySet))
					Expect(result.Sets[1]).To(BeIdenticalTo(kittenSet))
				})
			})

			Context("When there is no lift with that id in the database", func() {
				var (
					result *liftdatamodel.Lift
				)

				BeforeEach(func() {
					result = subject.GetById(1111111)
				})

				It("returns nil", func() {
					Expect(result).To(BeNil())
				})
			})
		})
	})

	Describe("Inserting lifts into the database", func() {
		var (
			createdId int
			err       error
		)

		BeforeEach(func() {
			setOne := &setdatamodel.Set{}
			setTwo := &setdatamodel.Set{}
			sets := []*setdatamodel.Set{setOne, setTwo}

			fakeSetRepository.InsertStub = func(set *setdatamodel.Set) (int, error) {
				if set == setOne {
					return 35, nil
				} else if set == setTwo {
					return 36, nil
				}

				return -1, nil
			}

			newLift := &liftdatamodel.Lift{
				Id:           -1,
				Name:         "turtle hang cleans",
				Workout:      25,
				DataTemplate: "time_in_seconds",
				Sets:         sets,
			}

			createdId, err = subject.Insert(newLift)
		})

		It("does not return an error", func() {
			Expect(err).To(BeNil())
		})

		It("returns the id of the created lift", func() {
			Expect(createdId).To(Equal(3))
		})

		It("creates a full lift in the database", func() {
			query := fmt.Sprintf("SELECT * FROM lifts WHERE id=%v", createdId)
			row := testConnection.QueryRow(query)

			var liftId int
			var name string
			var workout int
			var dataTemplate string
			var setIds sqlhelpers.IntSlice

			err := row.Scan(&liftId, &name, &workout, &dataTemplate, &setIds)
			Expect(err).To(BeNil())

			Expect(liftId).To(Equal(3))
			Expect(name).To(Equal("turtle hang cleans"))
			Expect(workout).To(Equal(25))
			Expect(dataTemplate).To(Equal("time_in_seconds"))

			Expect(len(setIds)).To(Equal(2))
			Expect(setIds[0]).To(Equal(35))
			Expect(setIds[1]).To(Equal(36))
		})
	})
})
