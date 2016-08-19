package translation_test

import (
	workoutdatamodel "github.com/jwfriese/workouttrackerapi/workouts/datamodel"
	workouttranslation "github.com/jwfriese/workouttrackerapi/workouts/translation"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("WorkoutsCreateRequestTranslator", func() {
	var (
		subject workouttranslation.WorkoutsCreateRequestTranslator
	)

	BeforeEach(func() {
		subject = workouttranslation.NewWorkoutsCreateRequestTranslator()
	})

	Describe("Translating request JSON into a workout", func() {
		var (
			result *workoutdatamodel.Workout
			err    error
		)

		Context("When the JSON is valid", func() {
			BeforeEach(func() {
				validJSON := []byte(`{"name":"turtle workout","timestamp":"2016-06-05T20:30:45-08:00"}`)
				result, err = subject.Translate(validJSON)
			})

			It("returns no error", func() {
				Expect(err).To(BeNil())
			})

			It("returns a valid workout model matching the request", func() {
				Expect(result).ToNot(BeNil())
				Expect(result.Name).To(Equal("turtle workout"))
				Expect(result.Timestamp).To(Equal("2016-06-05T20:30:45-08:00"))
				Expect(len(result.Lifts)).To(Equal(0))
			})
		})

		Context("Invalid JSON", func() {
			Context("Missing 'name'", func() {
				BeforeEach(func() {
					invalidJSON := []byte(`{"timestamp":"2016-06-05T20:30:45-08:00"}`)

					result, err = subject.Translate(invalidJSON)
				})

				It("returns an error", func() {
					Expect(result).To(BeNil())
					Expect(err.Error()).To(Equal("Missing required 'name' field from workout JSON"))
				})
			})

			Context("Missing 'timestamp'", func() {
				BeforeEach(func() {
					invalidJSON := []byte(`{"name":"turtle workout"}`)

					result, err = subject.Translate(invalidJSON)
				})

				It("returns an error", func() {
					Expect(result).To(BeNil())
					Expect(err.Error()).To(Equal("Missing required 'timestamp' field from workout JSON"))
				})
			})
		})
	})
})
