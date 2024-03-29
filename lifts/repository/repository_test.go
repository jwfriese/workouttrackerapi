package repository_test

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	liftdatamodel "github.com/jwfriese/workouttrackerapi/lifts/datamodel"
	liftrepository "github.com/jwfriese/workouttrackerapi/lifts/repository"
	setdatamodel "github.com/jwfriese/workouttrackerapi/lifts/sets/datamodel"
	setrepository "github.com/jwfriese/workouttrackerapi/lifts/sets/repository"
	setrepositoryfakes "github.com/jwfriese/workouttrackerapi/lifts/sets/repository/repositoryfakes"
	"github.com/jwfriese/workouttrackerapi/sqlhelpers"
	_ "github.com/lib/pq"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func populateDatabase(openConnection *sql.DB) {
	_, err := openConnection.Exec("INSERT INTO workouts (lifts) VALUES ('{}')")
	if err != nil {
		log.Fatal(err)
	}

	_, err = openConnection.Exec("INSERT INTO workouts (lifts) VALUES ('{}')")
	if err != nil {
		log.Fatal(err)
	}

	_, err = openConnection.Exec("INSERT INTO lifts (name,workout,data_template,sets) VALUES ('turtle lift',1,'weight/reps','{1,2}')")
	if err != nil {
		log.Fatal(err)
	}

	_, err = openConnection.Exec("INSERT INTO lifts (name,workout,data_template,sets) VALUES ('turtle jump',2,'height/reps','{3,4}')")
	if err != nil {
		log.Fatal(err)
	}
}

func clearDatabase(openConnection *sql.DB) {
	_, err := openConnection.Exec("TRUNCATE lifts")
	if err != nil {
		log.Fatal(err)
	}

	_, err = openConnection.Exec("ALTER SEQUENCE lifts_id_seq RESTART WITH 1")
	if err != nil {
		log.Fatal(err)
	}

	_, err = openConnection.Exec("TRUNCATE workouts")
	if err != nil {
		log.Fatal(err)
	}

	_, err = openConnection.Exec("ALTER SEQUENCE workouts_id_seq RESTART WITH 1")
	if err != nil {
		log.Fatal(err)
	}
}

var _ = Describe("LiftRepository", func() {
	var (
		subject           liftrepository.LiftRepository
		testConnection    *sql.DB
		fakeSetRepository *setrepositoryfakes.FakeSetRepository
	)

	BeforeEach(func() {
		var err error
		testConnection, err = sql.Open("postgres", "user=postgres dbname=workout_tracker_unit_test sslmode=disable")
		if err != nil {
			log.Fatal(err)
		}

		populateDatabase(testConnection)

		fakeSetRepository = new(setrepositoryfakes.FakeSetRepository)
		subject = liftrepository.NewLiftRepository(testConnection, fakeSetRepository)
	})

	AfterEach(func() {
		clearDatabase(testConnection)
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

			fakeSetRepository.GetByIdStub = func(id int) (*setdatamodel.Set, error) {
				if id == 1 {
					return turtleSet, nil
				} else if id == 2 {
					return crabSet, nil
				} else if id == 3 {
					return puppySet, nil
				} else if id == 4 {
					return kittenSet, nil
				}

				return nil, nil
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
			var (
				result *liftdatamodel.Lift
				err    error
			)

			Context("When there exists a lift with that ID in the database", func() {
				BeforeEach(func() {
					result, err = subject.GetById(2)
				})

				It("returns no error", func() {
					Expect(err).To(BeNil())
				})

				It("retrieves the requested lift from the database", func() {
					Expect(result).ToNot(BeNil())

					Expect(result.Id).To(Equal(2))
					Expect(result.Name).To(Equal("turtle jump"))
					Expect(result.Workout).To(Equal(2))
					Expect(result.DataTemplate).To(Equal("height/reps"))

					Expect(len(result.Sets)).To(Equal(2))
					Expect(result.Sets[0]).To(BeIdenticalTo(puppySet))
					Expect(result.Sets[1]).To(BeIdenticalTo(kittenSet))
				})
			})

			Context("When there is no lift with that id in the database", func() {
				BeforeEach(func() {
					result, err = subject.GetById(1111111)
				})

				It("returns a descriptive error", func() {
					Expect(err).ToNot(BeNil())
					Expect(err.Error()).To(Equal("Lift with id=1111111 does not exist"))
				})

				It("returns no result", func() {
					Expect(result).To(BeNil())
				})
			})

			Context("When the set repository errors during a fetch", func() {
				BeforeEach(func() {
					fakeSetRepository.GetByIdStub = func(id int) (*setdatamodel.Set, error) {
						errString := fmt.Sprintf("Error fetching set (id=%v)", id)
						return nil, errors.New(errString)
					}

					result, err = subject.GetById(2)
				})

				It("returns no lift", func() {
					Expect(result).To(BeNil())
				})

				It("returns a descriptive error", func() {
					Expect(err).ToNot(BeNil())
					Expect(err.Error()).To(Equal("Error fetching lift (id=2): Error fetching set (id=3)"))
				})
			})
		})
	})

	Describe("Inserting lifts into the database", func() {
		var (
			createdId int
			err       error
			setOne    *setdatamodel.Set
			setTwo    *setdatamodel.Set
			sets      []*setdatamodel.Set
		)

		BeforeEach(func() {
			setOne = &setdatamodel.Set{}
			setTwo = &setdatamodel.Set{}
			sets = []*setdatamodel.Set{setOne, setTwo}
		})

		Describe("When insert input is valid", func() {
			BeforeEach(func() {
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
					Workout:      2,
					DataTemplate: "timeInSeconds",
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
				Expect(workout).To(Equal(2))
				Expect(dataTemplate).To(Equal("timeInSeconds"))

				Expect(len(setIds)).To(Equal(2))
				Expect(setIds[0]).To(Equal(35))
				Expect(setIds[1]).To(Equal(36))
			})

			It("passes the sets along to the set repository to be inserted in there", func() {
				Expect(fakeSetRepository.InsertArgsForCall(0)).To(BeIdenticalTo(setOne))
				Expect(fakeSetRepository.InsertArgsForCall(1)).To(BeIdenticalTo(setTwo))
			})

			It("adds the new lift's id to the sets that are inserted", func() {
				insertedSetOne := fakeSetRepository.InsertArgsForCall(0)
				insertedSetTwo := fakeSetRepository.InsertArgsForCall(1)

				Expect(insertedSetOne.Lift).To(Equal(3))
				Expect(insertedSetTwo.Lift).To(Equal(3))
			})

			It("associates the specified workout with the new lift in the database", func() {
				workoutRow := testConnection.QueryRow("SELECT lifts FROM workouts WHERE id=2")
				var liftIds sqlhelpers.IntSlice
				workoutScanErr := workoutRow.Scan(&liftIds)

				Expect(workoutScanErr).To(BeNil())
				Expect(liftIds).To(BeEquivalentTo([]int{3}))
			})
		})

		Describe("Errors cases", func() {
			Context("When inserting a set errors", func() {
				BeforeEach(func() {
					fakeSetRepository.InsertStub = func(set *setdatamodel.Set) (int, error) {
						return -1, errors.New("Error inserting set")
					}

					newLift := &liftdatamodel.Lift{
						Id:           -1,
						Name:         "turtle hang cleans",
						Workout:      2,
						DataTemplate: "timeInSeconds",
						Sets:         sets,
					}

					createdId, err = subject.Insert(newLift)
				})

				It("returns an invalid id", func() {
					Expect(createdId).To(Equal(-1))
				})

				It("returns a descriptive error", func() {
					Expect(err).ToNot(BeNil())
					Expect(err.Error()).To(Equal("Error inserting lift: Error inserting set"))
				})
			})

			Context("When the specified workout does not exist", func() {
				BeforeEach(func() {
					newLift := &liftdatamodel.Lift{
						Id:           -1,
						Name:         "turtle hang cleans",
						Workout:      -1,
						DataTemplate: "timeInSeconds",
						Sets:         sets,
					}

					createdId, err = subject.Insert(newLift)
				})

				It("returns an invalid id", func() {
					Expect(createdId).To(Equal(-1))
				})

				It("returns a descriptive error", func() {
					Expect(err).ToNot(BeNil())
					Expect(err.Error()).To(Equal("Error inserting lift: Workout with id=-1 does not exist"))
				})
			})
		})
	})

	Describe("Deleting lifts from the database", func() {
		var (
			err error
		)

		Context("When a lift with the given id exists in the database", func() {
			Context("When deleting the lift's sets returns no errors", func() {
				BeforeEach(func() {
					fakeSetRepository.DeleteStub = func(id int) error {
						return nil
					}

					err = subject.Delete(1)
				})

				It("returns no error", func() {
					Expect(err).To(BeNil())
				})

				It("deletes the lift from the database", func() {
					row := testConnection.QueryRow("SELECT id FROM lifts WHERE id=1")
					var liftId int
					readErr := row.Scan(&liftId)
					Expect(readErr).ToNot(BeNil())
					Expect(readErr).To(Equal(sql.ErrNoRows))
				})

				It("tells the set repository to delete all the lift's sets", func() {
					Expect(fakeSetRepository.DeleteCallCount()).To(Equal(2))
					Expect(fakeSetRepository.DeleteArgsForCall(0)).To(Equal(1))
					Expect(fakeSetRepository.DeleteArgsForCall(1)).To(Equal(2))
				})
			})

			Context("When the lift's sets do not exist in the database", func() {
				BeforeEach(func() {
					fakeSetRepository.DeleteStub = func(id int) error {
						return setrepository.ErrDoesNotExist
					}

					err = subject.Delete(1)
				})

				It("returns no error", func() {
					Expect(err).To(BeNil())
				})

				It("deletes the lift from the database", func() {
					row := testConnection.QueryRow("SELECT id FROM lifts WHERE id=1")
					var liftId int
					readErr := row.Scan(&liftId)
					Expect(readErr).ToNot(BeNil())
					Expect(readErr).To(Equal(sql.ErrNoRows))
				})

				It("tells the set repository to delete all the lift's sets", func() {
					Expect(fakeSetRepository.DeleteCallCount()).To(Equal(2))
					Expect(fakeSetRepository.DeleteArgsForCall(0)).To(Equal(1))
					Expect(fakeSetRepository.DeleteArgsForCall(1)).To(Equal(2))
				})
			})

			Context("When any of the lift's sets cannot be deleted", func() {
				BeforeEach(func() {
					fakeSetRepository.DeleteStub = func(id int) error {
						return errors.New("error")
					}

					err = subject.Delete(1)
				})

				It("returns an error describing the failure", func() {
					Expect(err).ToNot(BeNil())
					Expect(err.Error()).To(Equal("Lift delete failed: Could not delete set with id=1"))
				})

				It("does not delete the lift from the database", func() {
					row := testConnection.QueryRow("SELECT id FROM lifts WHERE id=1")
					var liftId int
					readErr := row.Scan(&liftId)
					Expect(readErr).To(BeNil())
					Expect(liftId).To(Equal(1))
				})
			})
		})

		Context("When no lift with the given id exists in the database", func() {
			BeforeEach(func() {
				err = subject.Delete(1123432)
			})

			It("returns DoesNotExist error", func() {
				Expect(err).To(Equal(liftrepository.ErrDoesNotExist))
			})
		})
	})
})
