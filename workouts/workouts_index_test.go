package workouts_test

import (
	"github.com/jwfriese/workouttrackerapi/httpfakes"
	"github.com/jwfriese/workouttrackerapi/workouts"
	"github.com/jwfriese/workouttrackerapi/workouts/datamodel"
	"github.com/jwfriese/workouttrackerapi/workouts/repository/repositoryfakes"
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("workouts", func() {
	Describe("its index handler", func() {
		var (
			handler               http.Handler
			fakeResponseWriter    *httpfakes.FakeResponseWriter
			fakeWorkoutRepository *repositoryfakes.FakeWorkoutRepository
			turtleWorkout         *datamodel.Workout
			crabWorkout           *datamodel.Workout
		)

		BeforeEach(func() {
			fakeWorkoutRepository = new(repositoryfakes.FakeWorkoutRepository)
			handler = workouts.WorkoutsHandler(fakeWorkoutRepository)
		})

		It("returns a handler", func() {
			Expect(handler).ToNot(BeNil())
		})

		Describe("JSON served from index", func() {
			BeforeEach(func() {
				turtleWorkout = &datamodel.Workout{
					Id:        1234,
					Timestamp: "turtle timestamp",
					Lifts:     []string{"turtle lift", "turtle press"},
					Name:      "turtle workout",
				}

				crabWorkout = &datamodel.Workout{
					Id:        2345,
					Timestamp: "crab timestamp",
					Lifts:     []string{"crab lift", "crab press"},
					Name:      "crab workout",
				}

				fakeResponseWriter = new(httpfakes.FakeResponseWriter)

				workouts := []*datamodel.Workout{turtleWorkout, crabWorkout}
				fakeWorkoutRepository.AllReturns(workouts)

				handler.ServeHTTP(fakeResponseWriter, nil)
			})

			It("writes all workouts as JSON to the response", func() {
				Expect(fakeResponseWriter.WriteArgsForCall(0)).To(Equal([]byte(`[{"Id":1234,"Timestamp":"turtle timestamp","Lifts":["turtle lift","turtle press"],"Name":"turtle workout"},{"Id":2345,"Timestamp":"crab timestamp","Lifts":["crab lift","crab press"],"Name":"crab workout"}]`)))
			})
		})
	})
})
