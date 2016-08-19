package workouts_test

import (
	"net/http"

	"github.com/jwfriese/workouttrackerapi/test/http/httpfakes"
	"github.com/jwfriese/workouttrackerapi/workouts"
	workoutdatamodel "github.com/jwfriese/workouttrackerapi/workouts/datamodel"
	"github.com/jwfriese/workouttrackerapi/workouts/repository/repositoryfakes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("/workouts", func() {
	Describe("GET", func() {
		var (
			handler               http.Handler
			fakeResponseWriter    *httpfakes.FakeResponseWriter
			fakeWorkoutRepository *repositoryfakes.FakeWorkoutRepository
			turtleWorkout         *workoutdatamodel.Workout
			crabWorkout           *workoutdatamodel.Workout
		)

		BeforeEach(func() {
			fakeWorkoutRepository = new(repositoryfakes.FakeWorkoutRepository)
			handler = workouts.WorkoutsIndexHandler(fakeWorkoutRepository)
		})

		Describe("JSON served from index", func() {
			BeforeEach(func() {
				turtleWorkout = &workoutdatamodel.Workout{
					Id:        1234,
					Timestamp: "turtle timestamp",
					Lifts:     []int{2, 3},
					Name:      "turtle workout",
				}

				crabWorkout = &workoutdatamodel.Workout{
					Id:        2345,
					Timestamp: "crab timestamp",
					Lifts:     []int{10, 20, 30},
					Name:      "crab workout",
				}

				fakeResponseWriter = new(httpfakes.FakeResponseWriter)

				workouts := []*workoutdatamodel.Workout{turtleWorkout, crabWorkout}
				fakeWorkoutRepository.AllReturns(workouts)

				handler.ServeHTTP(fakeResponseWriter, nil)
			})

			It("writes all workouts as JSON to the response", func() {
				Expect(fakeResponseWriter.WriteArgsForCall(0)).To(MatchJSON([]byte(`[{"id":1234,"timestamp":"turtle timestamp","lifts":[2,3],"name":"turtle workout"},{"id":2345,"timestamp":"crab timestamp","lifts":[10,20,30],"name":"crab workout"}]`)))
			})
		})
	})
})
