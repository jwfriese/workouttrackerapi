package repository_test

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"

	liftdatamodel "github.com/jwfriese/workouttrackerapi/lifts/datamodel"
	"github.com/jwfriese/workouttrackerapi/lifts/repository"
	setdatamodel "github.com/jwfriese/workouttrackerapi/lifts/sets/datamodel"
	"github.com/jwfriese/workouttrackerapi/lifts/sets/repository/repositoryfakes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("LiftRepository", func() {
	var (
		subject           repository.LiftRepository
		testConnection    *sql.DB
		fakeSetRepository *repositoryfakes.FakeSetRepository
	)

	BeforeEach(func() {
		var err error
		testConnection, err = sql.Open("postgres", "user=postgres dbname=workout_tracker_test sslmode=disable")
		if err != nil {
			log.Fatal(err)
		}

		fakeSetRepository = new(repositoryfakes.FakeSetRepository)
		subject = repository.NewLiftRepository(testConnection, fakeSetRepository)
	})

	Describe("Getting all the lifts from DB", func() {
		var (
			lifts     []*liftdatamodel.Lift
			turtleSet *setdatamodel.Set
			crabSet   *setdatamodel.Set
			puppySet  *setdatamodel.Set
		)

		BeforeEach(func() {
			turtleSet = &setdatamodel.Set{
				Id:            777,
				DataTemplate:  "weight/reps",
				Lift:          4,
				Weight:        nil,
				Height:        nil,
				TimeInSeconds: nil,
				Reps:          nil,
			}

			crabSet = &setdatamodel.Set{
				Id:            888,
				DataTemplate:  "weight/reps",
				Lift:          4,
				Weight:        nil,
				Height:        nil,
				TimeInSeconds: nil,
				Reps:          nil,
			}

			puppySet = &setdatamodel.Set{
				Id:            999,
				DataTemplate:  "weight/reps",
				Lift:          4,
				Weight:        nil,
				Height:        nil,
				TimeInSeconds: nil,
				Reps:          nil,
			}

			fakeSetRepository.GetByIdStub = func(id uint) *setdatamodel.Set {
				if id == 10 {
					return turtleSet
				} else if id == 11 {
					return crabSet
				} else if id == 12 {
					return puppySet
				}

				return nil
			}

			lifts = subject.All()
		})

		It("should get all lifts from the database", func() {
			Expect(len(lifts)).To(Equal(12))
			Expect(lifts[0].Id).To(Equal(1))
			Expect(lifts[0].Name).To(Equal("turtle lift"))
			Expect(lifts[3].Id).To(Equal(4))

			Expect(len(lifts[3].Sets)).To(Equal(3))
			Expect(lifts[3].Sets[0].Id).To(Equal(777))
			Expect(lifts[3].Sets[1].Id).To(Equal(888))
			Expect(lifts[3].Sets[2].Id).To(Equal(999))
		})
	})
})
