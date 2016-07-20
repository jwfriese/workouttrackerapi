package translation_test

import (
	setdatamodel "github.com/jwfriese/workouttrackerapi/lifts/sets/datamodel"
	"github.com/jwfriese/workouttrackerapi/lifts/sets/translation"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("SetsCreateRequestTranslator", func() {
	var (
		subject translation.SetsCreateRequestTranslator
	)

	BeforeEach(func() {
		subject = translation.NewSetsCreateRequestTranslator()
	})

	Describe("Translating a set request", func() {
		var (
			result *setdatamodel.Set
			err    error
		)

		Context("When request JSON is valid", func() {
			Context("When dataTemplate is 'weight/reps'", func() {
				BeforeEach(func() {
					validJSON := []byte(`{"dataTemplate":"weight/reps","weight":135.0,"reps":12}`)
					result, err = subject.Translate(validJSON)
				})

				It("translates the JSON into a set", func() {
					Expect(err).To(BeNil())
					Expect(result).ToNot(BeNil())
					Expect(result.DataTemplate).To(Equal("weight/reps"))
					Expect(result.Weight).ToNot(BeNil())
					Expect(*(result.Weight)).To(BeEquivalentTo(135.0))
					Expect(result.Reps).ToNot(BeNil())
					Expect(*(result.Reps)).To(Equal(12))
				})
			})

			Context("When dataTemplate is 'height/reps'", func() {
				BeforeEach(func() {
					validJSON := []byte(`{"dataTemplate":"height/reps","height":24.0, "reps":8}`)
					result, err = subject.Translate(validJSON)
				})

				It("translates the JSON into a set", func() {
					Expect(err).To(BeNil())
					Expect(result).ToNot(BeNil())
					Expect(result.DataTemplate).To(Equal("height/reps"))
					Expect(result.Height).ToNot(BeNil())
					Expect(*(result.Height)).To(BeEquivalentTo(24))
					Expect(result.Reps).ToNot(BeNil())
					Expect(*(result.Reps)).To(Equal(8))
				})
			})

			Context("When dataTemplate is 'timeInSeconds'", func() {
				BeforeEach(func() {
					validJSON := []byte(`{"dataTemplate":"timeInSeconds","timeInSeconds":65.0}`)
					result, err = subject.Translate(validJSON)
				})

				It("translates the JSON into a set", func() {
					Expect(err).To(BeNil())
					Expect(result).ToNot(BeNil())
					Expect(result.DataTemplate).To(Equal("timeInSeconds"))
					Expect(result.TimeInSeconds).ToNot(BeNil())
					Expect(*(result.TimeInSeconds)).To(BeEquivalentTo(65.0))
				})
			})

			Context("When dataTemplate is 'weight/timeInSeconds'", func() {
				BeforeEach(func() {
					validJSON := []byte(`{"dataTemplate":"weight/timeInSeconds","timeInSeconds":65.0,"weight":55.0}`)
					result, err = subject.Translate(validJSON)
				})

				It("translates the JSON into a set", func() {
					Expect(err).To(BeNil())
					Expect(result).ToNot(BeNil())
					Expect(result.DataTemplate).To(Equal("weight/timeInSeconds"))
					Expect(result.Weight).ToNot(BeNil())
					Expect(*(result.Weight)).To(BeEquivalentTo(55.0))
					Expect(result.TimeInSeconds).ToNot(BeNil())
					Expect(*(result.TimeInSeconds)).To(BeEquivalentTo(65.0))
				})
			})
		})

		Context("When request JSON is invalid", func() {
			Context("When 'dataTemplate' is missing", func() {
				BeforeEach(func() {
					invalidJSON := []byte(`{"everythingElseInJSON":"does not matter"}`)
					result, err = subject.Translate(invalidJSON)
				})

				It("returns a descriptive error", func() {
					Expect(result).To(BeNil())
					Expect(err).ToNot(BeNil())
					Expect(err.Error()).To(Equal("Missing required 'dataTemplate' field in set request"))
				})
			})

			Context("When 'dataTemplate' is some unrecognized value", func() {
				BeforeEach(func() {
					invalidJSON := []byte(`{"dataTemplate":"turtles/are/bananas"}`)
					result, err = subject.Translate(invalidJSON)
				})

				It("returns a descriptive error", func() {
					Expect(result).To(BeNil())
					Expect(err).ToNot(BeNil())
					Expect(err.Error()).To(Equal("Unrecognized data template 'turtles/are/bananas'"))
				})
			})

			Context("When 'dataTemplate' is 'weight/reps'", func() {
				Context("When 'weight' is missing", func() {
					BeforeEach(func() {

						invalidJSON := []byte(`{"dataTemplate":"weight/reps","reps":10}`)
						result, err = subject.Translate(invalidJSON)
					})

					It("returns a descriptive error", func() {
						Expect(result).To(BeNil())
						Expect(err).ToNot(BeNil())
						Expect(err.Error()).To(Equal("Missing required 'weight' field in request for 'weight/reps' set"))
					})
				})

				Context("When 'reps' is missing", func() {
					BeforeEach(func() {

						invalidJSON := []byte(`{"dataTemplate":"weight/reps","weight":100.0}`)
						result, err = subject.Translate(invalidJSON)
					})

					It("returns a descriptive error", func() {
						Expect(result).To(BeNil())
						Expect(err).ToNot(BeNil())
						Expect(err.Error()).To(Equal("Missing required 'reps' field in request for 'weight/reps' set"))
					})
				})
			})

			Context("When 'dataTemplate' is 'height/reps'", func() {
				Context("When 'height' is missing", func() {
					BeforeEach(func() {

						invalidJSON := []byte(`{"dataTemplate":"height/reps","reps":10}`)
						result, err = subject.Translate(invalidJSON)
					})

					It("returns a descriptive error", func() {
						Expect(result).To(BeNil())
						Expect(err).ToNot(BeNil())
						Expect(err.Error()).To(Equal("Missing required 'height' field in request for 'height/reps' set"))
					})
				})

				Context("When 'reps' is missing", func() {
					BeforeEach(func() {

						invalidJSON := []byte(`{"dataTemplate":"height/reps","height":100.0}`)
						result, err = subject.Translate(invalidJSON)
					})

					It("returns a descriptive error", func() {
						Expect(result).To(BeNil())
						Expect(err).ToNot(BeNil())
						Expect(err.Error()).To(Equal("Missing required 'reps' field in request for 'height/reps' set"))
					})
				})
			})

			Context("When 'dataTemplate' is 'timeInSeconds'", func() {
				Context("When 'timeInSeconds' is missing", func() {
					BeforeEach(func() {

						invalidJSON := []byte(`{"dataTemplate":"timeInSeconds"}`)
						result, err = subject.Translate(invalidJSON)
					})

					It("returns a descriptive error", func() {
						Expect(result).To(BeNil())
						Expect(err).ToNot(BeNil())
						Expect(err.Error()).To(Equal("Missing required 'timeInSeconds' field in request for 'timeInSeconds' set"))
					})
				})
			})

			Context("When 'dataTemplate' is 'weight/timeInSeconds'", func() {
				Context("When 'weight' is missing", func() {
					BeforeEach(func() {

						invalidJSON := []byte(`{"dataTemplate":"weight/timeInSeconds","timeInSeconds":45.0}`)
						result, err = subject.Translate(invalidJSON)
					})

					It("returns a descriptive error", func() {
						Expect(result).To(BeNil())
						Expect(err).ToNot(BeNil())
						Expect(err.Error()).To(Equal("Missing required 'weight' field in request for 'weight/timeInSeconds' set"))
					})
				})

				Context("When 'timeInSeconds' is missing", func() {
					BeforeEach(func() {

						invalidJSON := []byte(`{"dataTemplate":"weight/timeInSeconds","weight":145.0}`)
						result, err = subject.Translate(invalidJSON)
					})

					It("returns a descriptive error", func() {
						Expect(result).To(BeNil())
						Expect(err).ToNot(BeNil())
						Expect(err.Error()).To(Equal("Missing required 'timeInSeconds' field in request for 'weight/timeInSeconds' set"))
					})
				})
			})
		})
	})
})
