package repository_test

import (
	"database/sql"
	"github.com/jwfriese/workouttracker/workouts/datamodel"
	"github.com/jwfriese/workouttracker/workouts/repository"
	_ "github.com/lib/pq"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"log"
)

var _ = Describe("WorkoutRepository", func() {
	var (
		subject        repository.WorkoutRepository
		testConnection *sql.DB
	)

	BeforeEach(func() {
		var err error
		testConnection, err = sql.Open("postgres", "user=postgres dbname=workout_tracker_test sslmode=disable")
		if err != nil {
			log.Fatal(err)
		}
		subject = repository.NewWorkoutRepository(testConnection)
	})

	Describe("Getting all workouts from the DB", func() {
		var (
			workouts []*datamodel.Workout
		)

		BeforeEach(func() {
			workouts = subject.All()
		})

		It("returns all workouts from the DB", func() {
			Expect(len(workouts)).To(Equal(5))
			firstWorkout := workouts[0]

			Expect(firstWorkout.Id).To(Equal(1))
			Expect(firstWorkout.Name).To(Equal("turtle one"))
			Expect(len(firstWorkout.Lifts)).To(Equal(3))
		})
	})

	Describe("Getting a single workout from the DB", func() {
		var (
			workout *datamodel.Workout
		)

		Context("When there exists a workout with that id in the DB", func() {
			BeforeEach(func() {
				workout = subject.GetById("2")
			})

			It("returns the workout from the DB with that id", func() {
				Expect(workout).ToNot(BeNil())
				Expect(workout.Id).To(Equal(2))
				Expect(workout.Name).To(Equal("turtle two"))
				Expect(len(workout.Lifts)).To(Equal(2))
				Expect(workout.Lifts[0]).To(Equal("turtle press"))
				Expect(workout.Lifts[1]).To(Equal("turtle cleans"))
			})
		})

		Context("When there is no workout with that id in the DB", func() {
			BeforeEach(func() {
				workout = subject.GetById("100")
			})

			It("returns nil", func() {
				Expect(workout).To(BeNil())
			})
		})
	})
})
