package repository_test

import (
	"database/sql"
	"fmt"
	"log"

	liftdatamodel "github.com/jwfriese/workouttrackerapi/lifts/datamodel"
	"github.com/jwfriese/workouttrackerapi/lifts/repository/repositoryfakes"
	"github.com/jwfriese/workouttrackerapi/sqlhelpers"
	workoutdatamodel "github.com/jwfriese/workouttrackerapi/workouts/datamodel"
	"github.com/jwfriese/workouttrackerapi/workouts/repository"
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
		subject            repository.WorkoutRepository
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
		subject = repository.NewWorkoutRepository(testConnection, fakeLiftRepository)
	})

	AfterEach(func() {
		clearWorkoutsDatabase(testConnection)
	})

	Describe("Getting workouts from the DB", func() {
		var (
			turtleLift *liftdatamodel.Lift
			crabLift   *liftdatamodel.Lift
			puppyLift  *liftdatamodel.Lift
		)

		BeforeEach(func() {
			turtleLift = &liftdatamodel.Lift{}
			crabLift = &liftdatamodel.Lift{}
			puppyLift = &liftdatamodel.Lift{}

			fakeLiftRepository.GetByIdStub = func(id int) *liftdatamodel.Lift {
				if id == 1 {
					return turtleLift
				} else if id == 2 {
					return crabLift
				} else if id == 3 {
					return puppyLift
				}

				return nil
			}
		})

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
				Expect(firstWorkout.Lifts[0]).To(BeIdenticalTo(turtleLift))
				Expect(firstWorkout.Lifts[1]).To(BeIdenticalTo(crabLift))

				secondWorkout := workouts[1]
				Expect(secondWorkout.Id).To(Equal(2))
				Expect(secondWorkout.Name).To(Equal("turtle two"))

				Expect(len(secondWorkout.Lifts)).To(Equal(1))
				Expect(secondWorkout.Lifts[0]).To(BeIdenticalTo(puppyLift))
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
					Expect(workout.Lifts[0]).To(BeIdenticalTo(puppyLift))
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
			turtleLift *liftdatamodel.Lift
			crabLift   *liftdatamodel.Lift
			createdId  int
			err        error
		)

		BeforeEach(func() {
			turtleLift = &liftdatamodel.Lift{}
			crabLift = &liftdatamodel.Lift{}
			lifts := []*liftdatamodel.Lift{turtleLift, crabLift}
			workout := &workoutdatamodel.Workout{
				Id:        -1,
				Timestamp: "1990-06-05T12:00:00-08:00",
				Name:      "name",
				Lifts:     lifts,
			}

			fakeLiftRepository.InsertStub = func(lift *liftdatamodel.Lift) (int, error) {
				if lift == turtleLift {
					return 10, nil
				} else if lift == crabLift {
					return 11, nil
				}

				return -1, nil
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
			Expect(len(liftIds)).To(Equal(2))
			Expect(liftIds[0]).To(Equal(10))
			Expect(liftIds[1]).To(Equal(11))
		})

		It("passes the lifts along to the lift repository to be inserted in there", func() {
			Expect(fakeLiftRepository.InsertArgsForCall(0)).To(BeIdenticalTo(turtleLift))
			Expect(fakeLiftRepository.InsertArgsForCall(1)).To(BeIdenticalTo(crabLift))
		})
	})
})
