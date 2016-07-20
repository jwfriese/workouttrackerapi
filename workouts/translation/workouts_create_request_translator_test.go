package translation_test

import (
	"bytes"
	"errors"
	"fmt"

	liftdatamodel "github.com/jwfriese/workouttrackerapi/lifts/datamodel"
	lifttranslationfakes "github.com/jwfriese/workouttrackerapi/lifts/translation/translationfakes"
	workoutdatamodel "github.com/jwfriese/workouttrackerapi/workouts/datamodel"
	workouttranslation "github.com/jwfriese/workouttrackerapi/workouts/translation"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("WorkoutsCreateRequestTranslator", func() {
	var (
		subject            workouttranslation.WorkoutsCreateRequestTranslator
		fakeLiftTranslator *lifttranslationfakes.FakeLiftsCreateRequestTranslator
	)

	BeforeEach(func() {
		fakeLiftTranslator = new(lifttranslationfakes.FakeLiftsCreateRequestTranslator)
		subject = workouttranslation.NewWorkoutsCreateRequestTranslator(fakeLiftTranslator)
	})

	Describe("Translating request JSON into a workout", func() {
		var (
			result *workoutdatamodel.Workout
			err    error
		)

		Context("When the JSON is valid", func() {
			var (
				liftOne *liftdatamodel.Lift
				liftTwo *liftdatamodel.Lift
			)

			BeforeEach(func() {
				liftOne = &liftdatamodel.Lift{}
				liftTwo = &liftdatamodel.Lift{}

				fakeLiftTranslator.TranslateStub = func(liftJSON []byte) (*liftdatamodel.Lift, error) {
					if bytes.Equal([]byte(`{"name":"liftOne"}`), liftJSON) {
						return liftOne, nil
					} else if bytes.Equal([]byte(`{"name":"liftTwo"}`), liftJSON) {
						return liftTwo, nil
					}

					errString := fmt.Sprintf("Invalid liftJSON (%q). Expected (%q) or (%q)", liftJSON, []byte(`{"name":"liftOne}"`), []byte(`{"name":"liftTwo"}`))
					return nil, errors.New(errString)
				}

				validJSON := []byte(`{"name":"turtle workout","timestamp":"2016-06-05T20:30:45-08:00","lifts":[{"name":"liftOne"},{"name":"liftTwo"}]}`)
				result, err = subject.Translate(validJSON)
			})

			It("returns no error", func() {
				Expect(err).To(BeNil())
			})

			It("returns a valid workout model matching the request", func() {
				Expect(result).ToNot(BeNil())
				Expect(result.Name).To(Equal("turtle workout"))
				Expect(result.Timestamp).To(Equal("2016-06-05T20:30:45-08:00"))
				Expect(len(result.Lifts)).To(Equal(2))
				Expect(result.Lifts[0]).To(BeIdenticalTo(liftOne))
				Expect(result.Lifts[1]).To(BeIdenticalTo(liftTwo))
			})
		})

		Context("Invalid JSON", func() {
			Context("Missing 'name'", func() {
				BeforeEach(func() {
					invalidJSON := []byte(`{"timestamp":"2016-06-05T20:30:45-08:00","lifts":["{liftOne}","{liftTwo}"]}`)

					result, err = subject.Translate(invalidJSON)
				})

				It("returns an error", func() {
					Expect(result).To(BeNil())
					Expect(err.Error()).To(Equal("Missing required 'name' field from workout JSON"))
				})
			})

			Context("Missing 'timestamp'", func() {
				BeforeEach(func() {
					invalidJSON := []byte(`{"name":"turtle workout","lifts":["{liftOne}","{liftTwo}"]}`)

					result, err = subject.Translate(invalidJSON)
				})

				It("returns an error", func() {
					Expect(result).To(BeNil())
					Expect(err.Error()).To(Equal("Missing required 'timestamp' field from workout JSON"))
				})
			})

			Context("Missing 'lifts'", func() {
				BeforeEach(func() {
					invalidJSON := []byte(`{"name":"turtle workout","timestamp":"2016-06-05T20:30:45-08:00"}`)

					result, err = subject.Translate(invalidJSON)
				})

				It("returns an error", func() {
					Expect(result).To(BeNil())
					Expect(err.Error()).To(Equal("Missing required 'lifts' field from workout JSON"))
				})
			})
		})
	})
})
