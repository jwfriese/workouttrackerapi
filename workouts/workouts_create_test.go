package workouts_test

import (
	"errors"

	"github.com/jwfriese/workouttrackerapi/httpfakes"
	liftrepositoryfakes "github.com/jwfriese/workouttrackerapi/lifts/repository/repositoryfakes"
	setrepositoryfakes "github.com/jwfriese/workouttrackerapi/lifts/sets/repository/repositoryfakes"
	"github.com/jwfriese/workouttrackerapi/workouts"
	workoutrepositoryfakes "github.com/jwfriese/workouttrackerapi/workouts/repository/repositoryfakes"
	"github.com/jwfriese/workouttrackerapi/workouts/validation/validationfakes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("POST /workouts", func() {
	Describe("Posting to workouts endpoint", func() {
		var (
			subject                  http.Handler
			fakeWorkoutRepository    *workoutrepositoryfakes.FakeWorkoutRepository
			fakeLiftRepository       *liftrepositoryfakes.FakeLiftRepository
			fakeValidator            *validationfakes.FakeWorkoutCreateRequestValidator
			fakeWorkoutReponseWriter *httpfakes.FakeReponseWriter
		)

		BeforeEach(func() {
			fakeResponseWriter = new(httpfakes.FakeResponseWriter)
			fakeWorkoutRepository = new(repositoryfakes.FakeWorkoutRepository)
			fakeValidator = new(validationfakes.FakeWorkoutCreateRequestValidator)
			subject = workouts.WorkoutsCreateHandler(fakeWorkoutRepository, fakeWorkoutCreateRequestValidator)
		})

		Describe("When request body fails validation", func() {
			BeforeEach(func() {
				err := errors.New("Validation error message")
				fakeValidator.ValidateReturns(err)
			})

			It("returns a 401", func() {
				Expect(fakeResponseWriter.WriteHeaderArgsForCall(0)).To(Equal(401))
			})

			It("returns an error message", func() {
				Expect(fakeResponseWriter.WriteArgsForCall(0)).To(Equal([]byte(`{"error":"Validation error message"}`)))
			})
		})
	})
})
