package repository_test

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	liftrepository "github.com/jwfriese/workouttrackerapi/lifts/repository"
	"github.com/jwfriese/workouttrackerapi/lifts/repository/repositoryfakes"
	"github.com/jwfriese/workouttrackerapi/sqlhelpers"
	workoutdatamodel "github.com/jwfriese/workouttrackerapi/workouts/datamodel"
	workoutrepository "github.com/jwfriese/workouttrackerapi/workouts/repository"
	_ "github.com/lib/pq"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func populateWorkoutsDatabase(openConnection *sql.DB) {
	_, err := openConnection.Exec("INSERT INTO workouts (name,timestamp,lifts) VALUES ('turtle one','2016-04-05T05:06:07-08:00','{1,2}')")
	if err != nil {
		log.Fatal(err)
	}

	_, err = openConnection.Exec("INSERT INTO workouts (name,timestamp,lifts) VALUES ('turtle two','2016-04-07T06:07:08-08:00','{3}')")
	if err != nil {
		log.Fatal(err)
	}
}

func clearWorkoutsDatabase(openConnection *sql.DB) {
	_, err := openConnection.Exec("TRUNCATE workouts")
	if err != nil {
		log.Fatal(err)
	}

	_, err = openConnection.Exec("ALTER SEQUENCE workouts_id_seq RESTART WITH 1")
	if err != nil {
		log.Fatal(err)
	}
}

var _ = Describe("WorkoutRepository", func() {
	var (
		subject            workoutrepository.WorkoutRepository
		testConnection     *sql.DB
		fakeLiftRepository *repositoryfakes.FakeLiftRepository
	)

	BeforeEach(func() {
		var err error
		testConnection, err = sql.Open("postgres", "user=postgres dbname=workout_tracker_unit_test sslmode=disable")
		if err != nil {
			log.Fatal(err)
		}

		populateWorkoutsDatabase(testConnection)
		fakeLiftRepository = new(repositoryfakes.FakeLiftRepository)
		subject = workoutrepository.NewWorkoutRepository(testConnection, fakeLiftRepository)
	})

	AfterEach(func() {
		clearWorkoutsDatabase(testConnection)
	})

	Describe("Getting workouts from the DB", func() {
		Describe("All", func() {
			var (
				workouts []*workoutdatamodel.Workout
			)

			BeforeEach(func() {
				workouts = subject.All()
			})

			It("returns all workouts from the DB", func() {
				Expect(len(workouts)).To(Equal(2))
				firstWorkout := workouts[0]

				Expect(firstWorkout.Id).To(Equal(1))
				Expect(firstWorkout.Name).To(Equal("turtle one"))

				Expect(len(firstWorkout.Lifts)).To(Equal(2))
				Expect(firstWorkout.Lifts[0]).To(Equal(1))
				Expect(firstWorkout.Lifts[1]).To(Equal(2))

				secondWorkout := workouts[1]
				Expect(secondWorkout.Id).To(Equal(2))
				Expect(secondWorkout.Name).To(Equal("turtle two"))

				Expect(len(secondWorkout.Lifts)).To(Equal(1))
				Expect(secondWorkout.Lifts[0]).To(Equal(3))
			})
		})

		Describe("A single workout", func() {
			var (
				workout *workoutdatamodel.Workout
				err     error
			)

			Context("When there exists a workout with that id in the DB", func() {
				BeforeEach(func() {
					workout, err = subject.GetById(2)
				})

				It("returns no error", func() {
					Expect(err).To(BeNil())
				})

				It("returns the workout from the DB with that id", func() {
					Expect(workout).ToNot(BeNil())
					Expect(workout.Id).To(Equal(2))
					Expect(workout.Name).To(Equal("turtle two"))
					Expect(len(workout.Lifts)).To(Equal(1))
					Expect(workout.Lifts[0]).To(Equal(3))
				})
			})

			Context("When there is no workout with that id in the DB", func() {
				BeforeEach(func() {
					workout, err = subject.GetById(100)
				})

				It("returns nil workout", func() {
					Expect(workout).To(BeNil())
				})

				It("returns a descriptive error", func() {
					Expect(err).ToNot(BeNil())
					Expect(err.Error()).To(Equal("Workout with id=100 does not exist"))
				})
			})
		})
	})

	Describe("Inserting a workout into the database", func() {
		var (
			createdId int
			err       error
		)

		Context("When workout input is valid", func() {
			BeforeEach(func() {
				workout := &workoutdatamodel.Workout{
					Id:        -1,
					Timestamp: "1990-06-05T12:00:00-08:00",
					Name:      "name",
					Lifts:     nil,
				}

				createdId, err = subject.Insert(workout)
			})

			It("does not return an error", func() {
				Expect(err).To(BeNil())
			})

			It("returns the id of the created workout", func() {
				Expect(createdId).To(Equal(3))
			})

			It("creates a full workout in the database", func() {
				query := fmt.Sprintf("SELECT * FROM workouts WHERE id=%v", createdId)
				row := testConnection.QueryRow(query)
				var id int
				var name string
				var timestamp string
				var liftIds sqlhelpers.IntSlice
				err := row.Scan(&id, &name, &timestamp, &liftIds)

				Expect(err).To(BeNil())
				Expect(name).To(Equal("name"))
				Expect(timestamp).To(SatisfyAny(Equal("1990-06-05T12:00:00-08:00"), Equal("1990-06-05T13:00:00-07:00")))
				Expect(len(liftIds)).To(Equal(0))
			})
		})

		Context("When the insert errors", func() {
			Context("When the timestamp is improperly formatted", func() {
				BeforeEach(func() {
					workout := &workoutdatamodel.Workout{
						Id:        -1,
						Timestamp: "not-a-timestamp",
						Name:      "name",
						Lifts:     nil,
					}

					createdId, err = subject.Insert(workout)
				})

				It("returns the error generated by the database", func() {
					Expect(err).ToNot(BeNil())
				})

				It("returns -1 for the id", func() {
					Expect(createdId).To(Equal(-1))
				})
			})

			Context("When workout input contains lift ids", func() {
				BeforeEach(func() {
					workout := &workoutdatamodel.Workout{
						Id:        -1,
						Timestamp: "1990-06-05T12:00:00-08:00",
						Name:      "name",
						Lifts:     []int{11, 21},
					}

					createdId, err = subject.Insert(workout)
				})

				It("returns an error", func() {
					Expect(err).ToNot(BeNil())
					Expect(err.Error()).To(Equal("Cannot include lift ids when inserting a workout (found [11 21])"))
				})

				It("returns -1 for the id", func() {
					Expect(createdId).To(Equal(-1))
				})
			})
		})
	})

	Describe("Deleting workouts", func() {
		var (
			err error
		)

		Context("When the workout exists in the database", func() {
			Context("When deleting workout lifts produces no errors", func() {
				BeforeEach(func() {
					fakeLiftRepository.DeleteStub = func(liftId int) error {
						return nil
					}

					err = subject.Delete(1)
				})

				It("returns no error", func() {
					Expect(err).To(BeNil())
				})

				It("deletes the workout with the given id from the database", func() {
					row := testConnection.QueryRow("SELECT id FROM workouts WHERE id=1")
					var id int
					readErr := row.Scan(&id)
					Expect(readErr).ToNot(BeNil())
					Expect(readErr).To(Equal(sql.ErrNoRows))
				})

				It("tells the lift repository to delete the associated lifts", func() {
					Expect(fakeLiftRepository.DeleteCallCount()).To(Equal(2))
					Expect(fakeLiftRepository.DeleteArgsForCall(0)).To(Equal(1))
					Expect(fakeLiftRepository.DeleteArgsForCall(1)).To(Equal(2))
				})
			})

			Context("When the workout's lifts do not exist", func() {
				BeforeEach(func() {
					fakeLiftRepository.DeleteStub = func(liftId int) error {
						return liftrepository.ErrDoesNotExist
					}

					err = subject.Delete(1)
				})

				It("returns no error", func() {
					Expect(err).To(BeNil())
				})

				It("deletes the workout with the given id from the database", func() {
					row := testConnection.QueryRow("SELECT id FROM workouts WHERE id=1")
					var id int
					readErr := row.Scan(&id)
					Expect(readErr).ToNot(BeNil())
					Expect(readErr).To(Equal(sql.ErrNoRows))
				})

				It("tells the lift repository to delete the associated lifts", func() {
					Expect(fakeLiftRepository.DeleteCallCount()).To(Equal(2))
					Expect(fakeLiftRepository.DeleteArgsForCall(0)).To(Equal(1))
					Expect(fakeLiftRepository.DeleteArgsForCall(1)).To(Equal(2))
				})
			})

			Context("When any of the workout's lifts fail to delete", func() {
				BeforeEach(func() {
					fakeLiftRepository.DeleteStub = func(id int) error {
						return errors.New("error")
					}

					err = subject.Delete(1)
				})

				It("returns a descriptive error", func() {
					Expect(err).ToNot(BeNil())
					Expect(err.Error()).To(Equal("Workout failed to delete: Lift with id=1 could not be deleted"))
				})

				It("does not delete the workout from the database", func() {
					row := testConnection.QueryRow("SELECT id FROM workouts WHERE id=1")
					var id int
					readErr := row.Scan(&id)
					Expect(readErr).To(BeNil())
					Expect(id).To(Equal(1))
				})
			})
		})

		Context("When the workout does not exist in the database", func() {
			BeforeEach(func() {
				err = subject.Delete(123456)
			})

			It("returns a DoesNotExist error", func() {
				Expect(err).ToNot(BeNil())
				Expect(err).To(Equal(workoutrepository.ErrDoesNotExist))
			})
		})
	})
})
