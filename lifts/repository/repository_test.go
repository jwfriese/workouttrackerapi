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
			turtleSet = &setdatamodel.Set{}
			crabSet = &setdatamodel.Set{}
			puppySet = &setdatamodel.Set{}

			fakeSetRepository.GetByIdStub = func(id int) *setdatamodel.Set {
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
			Expect(lifts[3].Sets[0]).To(BeIdenticalTo(turtleSet))
			Expect(lifts[3].Sets[1]).To(BeIdenticalTo(crabSet))
			Expect(lifts[3].Sets[2]).To(BeIdenticalTo(puppySet))
		})
	})

	Describe("Getting a single lift by ID from the database", func() {
		Context("When there exists a lift with that ID in the database", func() {
			var (
				result *liftdatamodel.Lift
				set    *setdatamodel.Set
			)

			BeforeEach(func() {
				set = &setdatamodel.Set{}
				fakeSetRepository.GetByIdStub = func(id int) *setdatamodel.Set {
					if id == 34 {
						return set
					}

					return nil
				}

				result = subject.GetById(12)
			})

			It("retrieves the requested lift from the database", func() {
				Expect(result).ToNot(BeNil())

				Expect(result.Id).To(Equal(12))
				Expect(result.Name).To(Equal("turtle hangs"))
				Expect(result.Workout).To(Equal(5))
				Expect(result.DataTemplate).To(Equal("time_in_seconds"))

				Expect(len(result.Sets)).To(Equal(1))
				Expect(result.Sets[0]).To(BeIdenticalTo(set))
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
