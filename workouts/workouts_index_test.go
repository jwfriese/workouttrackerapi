package workouts_test

import (
	"net/http"

	"github.com/jwfriese/workouttrackerapi/httpfakes"
	liftdatamodel "github.com/jwfriese/workouttrackerapi/lifts/datamodel"
	setdatamodel "github.com/jwfriese/workouttrackerapi/lifts/sets/datamodel"
	"github.com/jwfriese/workouttrackerapi/workouts"
	workoutdatamodel "github.com/jwfriese/workouttrackerapi/workouts/datamodel"
	"github.com/jwfriese/workouttrackerapi/workouts/repository/repositoryfakes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("workouts", func() {
	Describe("its index handler", func() {
		var (
			handler               http.Handler
			fakeResponseWriter    *httpfakes.FakeResponseWriter
			fakeWorkoutRepository *repositoryfakes.FakeWorkoutRepository
			turtleWorkout         *workoutdatamodel.Workout
			crabWorkout           *workoutdatamodel.Workout
		)

		BeforeEach(func() {
			fakeWorkoutRepository = new(repositoryfakes.FakeWorkoutRepository)
			handler = workouts.WorkoutsHandler(fakeWorkoutRepository)
		})

		It("returns a handler", func() {
			Expect(handler).ToNot(BeNil())
		})

		Describe("JSON served from index", func() {
			var (
				liftOne *liftdatamodel.Lift
				liftTwo *liftdatamodel.Lift
			)

			BeforeEach(func() {
				liftOne = &liftdatamodel.Lift{
					Id:           60,
					Name:         "lift one",
					DataTemplate: "dt1",
					Workout:      1234,
					Sets:         []*setdatamodel.Set{},
				}

				liftTwo = &liftdatamodel.Lift{
					Id:           61,
					Name:         "lift two",
					DataTemplate: "dt2",
					Workout:      1234,
					Sets:         []*setdatamodel.Set{},
				}

				turtleWorkout = &workoutdatamodel.Workout{
					Id:        1234,
					Timestamp: "turtle timestamp",
					Lifts:     []*liftdatamodel.Lift{liftOne, liftTwo},
					Name:      "turtle workout",
				}

				crabWorkout = &workoutdatamodel.Workout{
					Id:        2345,
					Timestamp: "crab timestamp",
					Lifts:     []*liftdatamodel.Lift{},
					Name:      "crab workout",
				}

				fakeResponseWriter = new(httpfakes.FakeResponseWriter)

				workouts := []*workoutdatamodel.Workout{turtleWorkout, crabWorkout}
				fakeWorkoutRepository.AllReturns(workouts)

				handler.ServeHTTP(fakeResponseWriter, nil)
			})

			It("writes all workouts as JSON to the response", func() {
				Expect(fakeResponseWriter.WriteArgsForCall(0)).To(Equal([]byte(`[{"id":1234,"timestamp":"turtle timestamp","lifts":[{"Id":60,"Name":"lift one","DataTemplate":"dt1","Workout":1234,"Sets":[]},{"Id":61,"Name":"lift two","DataTemplate":"dt2","Workout":1234,"Sets":[]}],"name":"turtle workout"},{"id":2345,"timestamp":"crab timestamp","lifts":[],"name":"crab workout"}]`)))
			})
		})
	})
})
