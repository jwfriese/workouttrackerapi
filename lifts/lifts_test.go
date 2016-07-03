package lifts_test

import (
	"net/http"

	"github.com/jwfriese/workouttrackerapi/httpfakes"
	"github.com/jwfriese/workouttrackerapi/lifts"
	liftdatamodel "github.com/jwfriese/workouttrackerapi/lifts/datamodel"
	"github.com/jwfriese/workouttrackerapi/lifts/repository/repositoryfakes"
	setdatamodel "github.com/jwfriese/workouttrackerapi/lifts/sets/datamodel"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("lifts", func() {
	Describe("its index handler", func() {
		var (
			handler            http.Handler
			fakeResponseWriter *httpfakes.FakeResponseWriter
			fakeLiftRepository *repositoryfakes.FakeLiftRepository
			turtlePressLift    *liftdatamodel.Lift
			turtleBoxJumpLift  *liftdatamodel.Lift
		)

		BeforeEach(func() {
			fakeLiftRepository = new(repositoryfakes.FakeLiftRepository)
			handler = lifts.LiftsHandler(fakeLiftRepository)
		})

		It("returns a handler", func() {
			Expect(handler).ToNot(BeNil())
		})

		Describe("JSON served from index", func() {
			BeforeEach(func() {
				fakeResponseWriter = new(httpfakes.FakeResponseWriter)

				turtlePressSets := make([]*setdatamodel.Set, 2)
				turtlePressSetOneWeight := new(float32)
				*turtlePressSetOneWeight = 123.4
				turtlePressSetOneReps := new(int)
				*turtlePressSetOneReps = 10
				turtlePressSets[0] = &setdatamodel.Set{
					Id:            1,
					DataTemplate:  "weight/reps",
					Weight:        turtlePressSetOneWeight,
					Height:        nil,
					TimeInSeconds: nil,
					Reps:          turtlePressSetOneReps,
				}

				turtlePressSetTwoWeight := new(float32)
				*turtlePressSetTwoWeight = 234.5
				turtlePressSetTwoReps := new(int)
				*turtlePressSetTwoReps = 6
				turtlePressSets[1] = &setdatamodel.Set{
					Id:            2,
					DataTemplate:  "weight/reps",
					Weight:        turtlePressSetTwoWeight,
					Height:        nil,
					TimeInSeconds: nil,
					Reps:          turtlePressSetTwoReps,
				}
				turtlePressLift = &liftdatamodel.Lift{
					Id:           1234,
					Name:         "turtle press",
					DataTemplate: "Weight/Reps",
					Workout:      1234567,
					Sets:         turtlePressSets,
				}

				turtleBoxJumpSets := make([]*setdatamodel.Set, 2)
				turtleBoxJumpSetOneHeight := new(float32)
				*turtleBoxJumpSetOneHeight = 36.0
				turtleBoxJumpSetOneReps := new(int)
				*turtleBoxJumpSetOneReps = 8
				turtleBoxJumpSets[0] = &setdatamodel.Set{
					Id:            3,
					DataTemplate:  "height/reps",
					Weight:        nil,
					Height:        turtleBoxJumpSetOneHeight,
					TimeInSeconds: nil,
					Reps:          turtleBoxJumpSetOneReps,
				}

				turtleBoxJumpSetTwoHeight := new(float32)
				*turtleBoxJumpSetTwoHeight = 42.0
				turtleBoxJumpSetTwoReps := new(int)
				*turtleBoxJumpSetTwoReps = 10
				turtleBoxJumpSets[1] = &setdatamodel.Set{
					Id:            4,
					DataTemplate:  "height/reps",
					Weight:        nil,
					Height:        turtleBoxJumpSetTwoHeight,
					TimeInSeconds: nil,
					Reps:          turtleBoxJumpSetTwoReps,
				}
				turtleBoxJumpLift = &liftdatamodel.Lift{
					Id:           2345,
					Name:         "turtle box jump",
					DataTemplate: "Height/Reps",
					Workout:      7654321,
					Sets:         turtleBoxJumpSets,
				}

				lifts := []*liftdatamodel.Lift{turtlePressLift, turtleBoxJumpLift}
				fakeLiftRepository.AllReturns(lifts)

				handler.ServeHTTP(fakeResponseWriter, nil)
			})

			It("writes all lift data to the response", func() {
				Expect(fakeResponseWriter.WriteArgsForCall(0)).To(Equal([]byte(`[{"id":1234,"name":"turtle press","dataTemplate":"Weight/Reps","workout":1234567,"sets":[{"weight":123.4,"reps":10},{"weight":234.5,"reps":6}]},{"id":2345,"name":"turtle box jump","dataTemplate":"Height/Reps","workout":7654321,"sets":[{"height":36,"reps":8},{"height":42,"reps":10}]}]`)))
			})
		})
	})
})
