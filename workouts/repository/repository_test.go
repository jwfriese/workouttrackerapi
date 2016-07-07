package repository_test

import (
	"database/sql"
	"log"

	liftdatamodel "github.com/jwfriese/workouttrackerapi/lifts/datamodel"
	"github.com/jwfriese/workouttrackerapi/lifts/repository/repositoryfakes"
	workoutdatamodel "github.com/jwfriese/workouttrackerapi/workouts/datamodel"
	"github.com/jwfriese/workouttrackerapi/workouts/repository"
	_ "github.com/lib/pq"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

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

		fakeLiftRepository = new(repositoryfakes.FakeLiftRepository)
		subject = repository.NewWorkoutRepository(testConnection, fakeLiftRepository)
	})

	Describe("Getting all workouts from the DB", func() {
		var (
			workouts   []*workoutdatamodel.Workout
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

			workouts = subject.All()
		})

		It("returns all workouts from the DB", func() {
			Expect(len(workouts)).To(Equal(5))
			firstWorkout := workouts[0]

			Expect(firstWorkout.Id).To(Equal(1))
			Expect(firstWorkout.Name).To(Equal("turtle one"))

			Expect(len(firstWorkout.Lifts)).To(Equal(3))
			Expect(firstWorkout.Lifts[0]).To(BeIdenticalTo(turtleLift))
			Expect(firstWorkout.Lifts[1]).To(BeIdenticalTo(crabLift))
			Expect(firstWorkout.Lifts[2]).To(BeIdenticalTo(puppyLift))
		})
	})

	Describe("Getting a single workout from the DB", func() {
		var (
			workout    *workoutdatamodel.Workout
			turtleLift *liftdatamodel.Lift
			crabLift   *liftdatamodel.Lift
		)

		Context("When there exists a workout with that id in the DB", func() {
			BeforeEach(func() {
				turtleLift = &liftdatamodel.Lift{}
				crabLift = &liftdatamodel.Lift{}

				fakeLiftRepository.GetByIdStub = func(id int) *liftdatamodel.Lift {
					if id == 4 {
						return turtleLift
					} else if id == 5 {
						return crabLift
					}

					return nil
				}
				workout = subject.GetById(2)
			})

			It("returns the workout from the DB with that id", func() {
				Expect(workout).ToNot(BeNil())
				Expect(workout.Id).To(Equal(2))
				Expect(workout.Name).To(Equal("turtle two"))
				Expect(len(workout.Lifts)).To(Equal(2))
				Expect(workout.Lifts[0]).To(BeIdenticalTo(turtleLift))
				Expect(workout.Lifts[1]).To(BeIdenticalTo(crabLift))
			})
		})

		Context("When there is no workout with that id in the DB", func() {
			BeforeEach(func() {
				workout = subject.GetById(100)
			})

			It("returns nil", func() {
				Expect(workout).To(BeNil())
			})
		})
	})
})
